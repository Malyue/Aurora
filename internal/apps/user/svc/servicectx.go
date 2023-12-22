package svc

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServerCtx struct {
	Ctx      context.Context
	Logger   *logrus.Logger
	DBClient *gorm.DB
	RedisCli *redis.Client
	Cache    map[string]Item
}

type Item struct {
}
