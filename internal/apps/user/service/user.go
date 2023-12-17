package service

import (
	userpb "Aurora/api/proto-go/user"
	"context"
	"sync"
)

var UserSvcOnce sync.Once
var UserSvcInstance *UserSvc

type UserSvc struct {
	userpb.UnimplementedUserServiceServer
}

func GetUserSvc() *UserSvc {
	UserSvcOnce.Do(func() {
		UserSvcInstance = &UserSvc{}
	})
	return UserSvcInstance
}

func (s *UserSvc) Hello(ctx context.Context, req *userpb.HelloRequest) (resp *userpb.HelloResponse, err error) {
	resp = &userpb.HelloResponse{
		Msg: "hello",
	}
	return
}
