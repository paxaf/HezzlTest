package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/nats-io/nats.go"
	"github.com/paxaf/HezzlTest/internal/entity"
	"github.com/paxaf/HezzlTest/internal/logger"
)

const (
	batchSize     = 1000
	flushInterval = 5 * time.Second
	maxRetries    = 5
)

type EventWithAck struct {
	Event *entity.Event
	Msg   *nats.Msg
}

type ClickHouseWorker struct {
	conn      driver.Conn
	nc        *nats.Conn
	js        nats.JetStreamContext
	batch     []EventWithAck
	batchLock sync.Mutex
	subject   string
	shutdown  chan struct{}
	wg        sync.WaitGroup
}

func NewClickHouseWorker(nc *nats.Conn, conn driver.Conn, subject string) (*ClickHouseWorker, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &ClickHouseWorker{
		conn:     conn,
		nc:       nc,
		js:       js,
		batch:    make([]EventWithAck, 0, batchSize*2),
		subject:  subject,
		shutdown: make(chan struct{}),
	}, nil
}

func (w *ClickHouseWorker) Start() {
	sub, err := w.js.Subscribe(w.subject, func(msg *nats.Msg) {
		event, err := w.unmarshalEvent(msg.Data)
		if err != nil {
			logger.Error("failed to unmarshal event", err)
			msg.Nak() // Negative acknowledgment
			return
		}

		w.AddEvent(event, msg)
	}, nats.Durable("clickhouse-worker"), nats.ManualAck())

	if err != nil {
		logger.Fatal("Failed to subscribe:", err)
	}

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer sub.Unsubscribe()
		w.batchProcessor()
	}()
}

func (w *ClickHouseWorker) unmarshalEvent(data []byte) (*entity.Event, error) {
	var base entity.BaseEvent
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	event := &entity.Event{BaseEvent: base}

	switch base.Entity {
	case "project":
		var payload entity.ProjectEventPayload
		if err := json.Unmarshal(data, &struct {
			Payload *entity.ProjectEventPayload `json:"payload"`
		}{Payload: &payload}); err != nil {
			return nil, err
		}
		event.Payload = payload
	case "good":
		var payload entity.GoodEventPayload
		if err := json.Unmarshal(data, &struct {
			Payload *entity.GoodEventPayload `json:"payload"`
		}{Payload: &payload}); err != nil {
			return nil, err
		}
		event.Payload = payload
	default:
		return nil, fmt.Errorf("unknown entity type: %s", base.Entity)
	}

	return event, nil
}

func (w *ClickHouseWorker) AddEvent(event *entity.Event, msg *nats.Msg) {
	w.batchLock.Lock()
	defer w.batchLock.Unlock()
	w.batch = append(w.batch, EventWithAck{Event: event, Msg: msg})
}

func (w *ClickHouseWorker) batchProcessor() {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.flushBatch()
		case <-w.shutdown:
			w.flushBatch()
			return
		}
	}
}

func (w *ClickHouseWorker) flushBatch() {
	w.batchLock.Lock()
	if len(w.batch) == 0 {
		w.batchLock.Unlock()
		return
	}

	batch := make([]EventWithAck, len(w.batch))
	copy(batch, w.batch)
	w.batch = w.batch[:0]
	w.batchLock.Unlock()

	successful, failed := w.insertBatchWithRetry(batch, maxRetries)
	w.ackMessages(successful)
	w.returnFailedToQueue(failed)
}

func (w *ClickHouseWorker) insertBatchWithRetry(batch []EventWithAck, maxRetries int) (successful, failed []EventWithAck) {
	for i := 0; i < maxRetries; i++ {
		successful, failed = w.insertBatch(batch)
		if len(failed) == 0 {
			return successful, nil
		}

		if i < maxRetries-1 {
			delay := time.Duration(i+1) * time.Second
			logger.Debug("Retrying failed events", "count", len(failed), "retry", i+1, "delay", delay)
			time.Sleep(delay)
			batch = failed
		}
	}
	return successful, failed
}

func (w *ClickHouseWorker) insertBatch(items []EventWithAck) (successful, failed []EventWithAck) {
	batch, err := w.conn.PrepareBatch(context.Background(), "INSERT INTO logs.events")
	if err != nil {
		logger.Error("failed to prepare batch", err)
		return nil, items // Все элементы не обработаны
	}

	for _, item := range items {
		event := item.Event
		err := w.appendEventToBatch(batch, event)
		if err != nil {
			logger.Error("failed to append event to batch", err, "entity_id", event.EntityID)
			failed = append(failed, item)
			continue
		}
		successful = append(successful, item)
	}

	if len(successful) == 0 {
		return nil, items // Нечего вставлять
	}

	if err := batch.Send(); err != nil {
		logger.Error("failed to send batch", err)
		return nil, items // Весь батч не обработан
	}

	log.Printf("Inserted %d events", len(successful))
	return successful, failed
}

func (w *ClickHouseWorker) appendEventToBatch(batch driver.Batch, event *entity.Event) error {
	switch event.Entity {
	case "project":
		payload := event.Payload.(entity.ProjectEventPayload)
		return batch.Append(
			event.Timestamp,
			string(event.Action),
			event.Entity,
			event.EntityID,
			event.ProjectID,
			payload.Name,
			nil,   // NULL for description
			0,     // Default priority
			false, // Default removed
			payload.CreatedAt,
		)
	case "good":
		payload := event.Payload.(entity.GoodEventPayload)
		return batch.Append(
			event.Timestamp,
			string(event.Action),
			event.Entity,
			event.EntityID,
			event.ProjectID,
			payload.Name,
			payload.Description,
			payload.Priority,
			payload.Removed,
			payload.CreatedAt,
		)
	default:
		return fmt.Errorf("unknown entity type: %s", event.Entity)
	}
}

func (w *ClickHouseWorker) ackMessages(items []EventWithAck) {
	for _, item := range items {
		if err := item.Msg.Ack(); err != nil {
			logger.Error("failed to ACK message", err)
		}
	}
}

func (w *ClickHouseWorker) returnFailedToQueue(items []EventWithAck) {
	for _, item := range items {
		data, err := json.Marshal(item.Event)
		if err != nil {
			logger.Error("failed to marshal event for requeue", err)
			continue
		}

		if _, err := w.js.Publish(w.subject, data); err != nil {
			logger.Error("failed to requeue message", err)
		}
	}
}

func (w *ClickHouseWorker) Close() {
	close(w.shutdown)
	w.wg.Wait()
	w.conn.Close()
}
