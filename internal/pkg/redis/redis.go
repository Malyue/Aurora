package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func NewRedis(cfg *Config) (*redis.Client, error) {
	if cfg == nil {
		return nil, errors.New("the config is nil")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return rdb, nil
}
