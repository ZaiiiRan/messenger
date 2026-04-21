package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func New(ctx context.Context, cfg settings.RedisSettings) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,

		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,

		PoolSize:        int(cfg.MaxPoolSize),
		MinIdleConns:    int(cfg.MinPoolSize),
		ConnMaxIdleTime: time.Duration(cfg.MaxConnIdleTime) * time.Second,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &RedisClient{client: rdb}, nil
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

func (r *RedisClient) Close() {
	if r.client != nil {
		r.client.Close()
	}
}
