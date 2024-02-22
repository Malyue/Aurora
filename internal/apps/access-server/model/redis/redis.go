package redis

import "github.com/go-redis/redis/v8"

type RedisClient struct {
	*redis.Client
}
