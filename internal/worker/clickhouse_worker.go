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

type ClickHouseWorker struct {
	conn      driver.Conn
	nc        *nats.Conn
	js        nats.JetStreamContext
	batch     []entity.Event
	batchLock sync.Mutex
	subject   string
}

func NewClickHouseWorker(nc *nats.Conn, conn driver.Conn, subject string) (*ClickHouseWorker, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &ClickHouseWorker{
		conn:    conn,
		nc:      nc,
		js:      js,
		batch:   make([]entity.Event, 0, batchSize*2),
		subject: subject,
	}, nil
}

func (w *ClickHouseWorker) Start(ctx context.Context) {
	_, err := w.js.Subscribe(w.subject, func(msg *nats.Msg) {
		var event entity.Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			logger.Error("failed to unmarshal event", err)
			return
		}

		w.AddEvent(event)
	}, nats.Durable("clickhouse-worker"), nats.ManualAck())

	if err != nil {
		logger.Fatal("Failed to subscribe:", err)
	}

	go w.batchProcessor(ctx)
}

func (w *ClickHouseWorker) AddEvent(event entity.Event) {
	w.batchLock.Lock()
	defer w.batchLock.Unlock()
	w.batch = append(w.batch, event)
}

func (w *ClickHouseWorker) batchProcessor(ctx context.Context) {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.flushBatch(ctx)
		case <-ctx.Done():
			w.flushBatch(ctx)
			return
		}
	}
}

func (w *ClickHouseWorker) flushBatch(ctx context.Context) {
	w.batchLock.Lock()
	if len(w.batch) == 0 {
		w.batchLock.Unlock()
		return
	}

	batch := make([]entity.Event, len(w.batch))
	copy(batch, w.batch)
	w.batch = w.batch[:0]
	w.batchLock.Unlock()

	for i := 0; i < maxRetries; i++ {
		if err := w.insertBatch(ctx, batch); err == nil {
			return
		}

		if i < maxRetries-1 {
			delay := time.Duration(i+1) * time.Second
			logger.Debug("Retrying batch insert in ", delay)
			time.Sleep(delay)
		}
	}

	logger.Error("Failed to insert batch after retries:", maxRetries)
}

func (w *ClickHouseWorker) insertBatch(ctx context.Context, events []entity.Event) error {
	batch, err := w.conn.PrepareBatch(ctx, "INSERT INTO logs.events")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, event := range events {
		if event.Entity == "project" {
			payload := event.Payload.(entity.ProjectEventPayload)
			err = batch.Append(
				event.Timestamp,
				string(event.Action),
				event.Entity,
				event.EntityID,
				event.ProjectID,
				payload.Name,
				"",
				0,
				false,
				payload.CreatedAt,
			)
		} else if event.Entity == "good" {
			payload := event.Payload.(entity.GoodEventPayload)
			err = batch.Append(
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
		}

		if err != nil {
			return fmt.Errorf("failed to append event: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	log.Printf("Successfully inserted %d events", len(events))
	return nil
}

func (ch *ClickHouseWorker) Close() {
	ch.conn.Close()
}
