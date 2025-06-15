package repository

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/paxaf/HezzlTest/config"
	"github.com/paxaf/HezzlTest/internal/logger"
)

type ClickHouse struct {
	conn driver.Conn
}

func NewClickHouse(cfg config.Clickhouse) (*ClickHouse, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.Address},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		DialTimeout: 10 * time.Second,
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	logger.Info("Connected to ClickHouse at", cfg.Address)
	return &ClickHouse{conn: conn}, nil
}

func (ch *ClickHouse) Migrate(ctx context.Context) error {
	return nil
}

func (ch *ClickHouse) Close() error {
	return ch.conn.Close()
}
