package svc

import (
	userpb "Aurora/api/proto-go/user"
	_redisx "Aurora/internal/apps/access-server/model/redis"
	"Aurora/internal/pkg/timingWheel"
	"context"
	"github.com/sirupsen/logrus"
)

type ServerCtx struct {
	Ctx      context.Context
	Logger   *logrus.Logger
	stopChan chan struct{}
	// grpc client
	UserServer  userpb.UserServiceClient
	RedisClient *_redisx.RedisClient
	TimingWheel *timingWheel.TimingWheel
}
