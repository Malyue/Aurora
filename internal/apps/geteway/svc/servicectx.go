package svc

import (
	userpb "Aurora/api/proto-go/user"
	"context"
	"github.com/sirupsen/logrus"
)

type ServerCtx struct {
	Ctx    context.Context
	Logger *logrus.Logger
	//Etcd   *etcd
	// Grpc Server
	UserServer userpb.UserServiceClient
}
