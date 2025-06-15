package events

import (
	"log"
	"time"

	"github.com/paxaf/HezzlTest/internal/entity"
	natsClient "github.com/paxaf/HezzlTest/internal/infrastructure/nats"
)

type Event struct {
	nats      *natsClient.NatsClient
	eventChan chan entity.Event
}

func New(nats *natsClient.NatsClient) *Event {
	nr := &Event{
		nats:      nats,
		eventChan: make(chan entity.Event, 1000),
	}
	go nr.eventProcessor()

	return nr
}

func (nr *Event) LogEvent(event entity.Event) {
	select {
	case nr.eventChan <- event:
	default:
		log.Printf("NATS event channel full, dropping event: %+v", event)
	}
}

func (nr *Event) eventProcessor() {
	const batchSize = 100
	batch := make([]entity.Event, 0, batchSize)
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case event := <-nr.eventChan:
			batch = append(batch, event)
			if len(batch) >= batchSize {
				nr.sendBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				nr.sendBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

func (nr *Event) sendBatch(events []entity.Event) {
	for _, event := range events {
		msg, err := event.Marshal()
		if err != nil {
			log.Printf("Failed to marshal event: %v", err)
			continue
		}

		subject := "db.events." + string(event.Action) + "." + event.Entity
		_, err = nr.nats.JS.PublishAsync(subject, msg)
		if err != nil {
			log.Printf("Failed to publish event: %v", err)
		}
	}
}

func (nr *Event) Close() {
	nr.nats.Close()
}
