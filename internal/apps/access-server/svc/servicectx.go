package svc

import (
	_redisx "Aurora/internal/apps/access-server/model/redis"
	"context"
	"github.com/sirupsen/logrus"
)

type ServerCtx struct {
	Ctx         context.Context
	Logger      *logrus.Logger
	stopChan    chan struct{}
	RedisClient *_redisx.RedisClient
}
