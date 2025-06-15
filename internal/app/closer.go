package app

import (
	"context"
	"fmt"
	"time"

	"github.com/paxaf/HezzlTest/internal/logger"
	"github.com/paxaf/HezzlTest/internal/repository/events"
	"github.com/paxaf/HezzlTest/internal/repository/postgres"
	redisClient "github.com/paxaf/HezzlTest/internal/repository/redis"
)

type closer struct {
	postgres *postgres.PgPool
	redis    *redisClient.RedisClient
	nats     *events.Event
}

func NewCloser(postgres *postgres.PgPool, redis *redisClient.RedisClient, nats *events.Event) *closer {
	return &closer{
		postgres: postgres,
		redis:    redis,
		nats:     nats,
	}
}

func (c *closer) Close(app *App) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if app.apiServer != nil {
		if err := app.apiServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("HTTP server shutdown failed: %w", err)
		}
	}

	c.postgres.Close()
	c.redis.Close()
	c.nats.Close()
	logger.Info("Database connections closed successfully")

	logger.Info("Application stopped gracefully")
	return nil
}
