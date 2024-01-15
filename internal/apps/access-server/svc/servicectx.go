package svc

import (
	"context"
	"github.com/sirupsen/logrus"
)

type ServerCtx struct {
	Ctx      context.Context
	Logger   *logrus.Logger
	stopChan chan struct{}
}
