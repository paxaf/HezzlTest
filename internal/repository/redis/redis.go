package redisClient

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/paxaf/HezzlTest/internal/entity"
)

const (
	ttl = 60 * time.Second
)

type RedisClient struct {
	client *redis.Client
}

func New(client *redis.Client) *RedisClient {
	return &RedisClient{client: client}
}

func (rc *RedisClient) RedisSetItem(key string, item interface{}) error {
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("redis failed marshal item: %w", err)
	}
	return rc.client.Set(key, data, ttl).Err()
}

func (rc *RedisClient) RedisGetItem(key string) (*entity.Goods, error) {
	data, err := rc.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed redis get item: %w", err)
	}
	var goods entity.Goods
	if err = json.Unmarshal([]byte(data), &goods); err != nil {
		return nil, fmt.Errorf("redis unmarshal error: %w", err)
	}
	return &goods, nil
}

func (rc *RedisClient) RedisGetItems(key string) ([]entity.Goods, error) {
	data, err := rc.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed redis get item: %w", err)
	}
	var goods []entity.Goods
	if err = json.Unmarshal([]byte(data), &goods); err != nil {
		return nil, fmt.Errorf("redis unmarshal error: %w", err)
	}
	return goods, nil
}

func (rc *RedisClient) CleanCache() error {
	return rc.client.FlushAll().Err()
}
