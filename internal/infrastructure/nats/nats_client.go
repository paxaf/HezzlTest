package natsClient

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/paxaf/HezzlTest/config"
)

type NatsClient struct {
	Conn *nats.Conn
	JS   nats.JetStreamContext
}

func New(cfg config.Nats) (*NatsClient, error) {
	nc, err := nats.Connect(cfg.Url,
		nats.Name("HezzlTest-App"),
		nats.Timeout(10*time.Second),
		nats.PingInterval(1*time.Minute),
		nats.MaxPingsOutstanding(5),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(-1),
	)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "DB_EVENTS",
		Subjects: []string{"db.events.>"},
		MaxAge:   24 * time.Hour,
	})
	if err != nil {
		log.Printf("Stream creation warning: %v", err)
	}

	log.Printf("Connected to NATS at: %s", cfg.Url)
	return &NatsClient{Conn: nc, JS: js}, nil
}

func (c *NatsClient) Close() {
	if c.Conn != nil && c.Conn.IsConnected() {
		c.Conn.Close()
	}
}
