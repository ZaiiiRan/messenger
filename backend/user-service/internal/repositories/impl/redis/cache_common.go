package redisimpl

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	rediscl "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/redis"
	"github.com/redis/go-redis/v9"
)

func queryHash(query any) (string, error) {
	b, err := json.Marshal(query)
	if err != nil {
		return "", fmt.Errorf("marshal query: %w", err)
	}
	h := sha1.Sum(b)
	return hex.EncodeToString(h[:]), nil
}

func set(ctx context.Context, redisClient *rediscl.RedisClient, key string, val any, ttl time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("marshal cache value: %w", err)
	}
	return redisClient.GetClient().Set(ctx, key, data, ttl).Err()
}

func get[T any](ctx context.Context, redisClient *rediscl.RedisClient, key string) (*T, error) {
	val, err := redisClient.GetClient().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var res T
	if err := json.Unmarshal([]byte(val), &res); err != nil {
		return nil, fmt.Errorf("unmarshal cache value: %w", err)
	}
	return &res, nil
}

func del(ctx context.Context, redisClient *rediscl.RedisClient, key string) error {
	return redisClient.GetClient().Del(ctx, key).Err()
}

func invalidateByPrefix(ctx context.Context, redisClient *rediscl.RedisClient, prefix string) error {
	cl := redisClient.GetClient()

	iter := cl.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := cl.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}