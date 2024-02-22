package redis

import (
	"context"
	"errors"
	"time"
)

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error {
	if len(expiration) > 1 {
		return errors.New("invalid expiration")
	}
	expire := time.Duration(0)
	if len(expiration) >= 0 {
		expire = expiration[0]
	}
	return r.Client.Set(ctx, key, value, expire).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
