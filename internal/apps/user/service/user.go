package service

import (
	userpb "Aurora/api/proto-go/user"
	"context"
	"sync"
)

var UserServerOnce sync.Once
var UserServerInstance *UserServer

type UserServer struct {
	userpb.UnimplementedUserServiceServer
}

var _ userpb.UserServiceServer = (*UserServer)(nil)

func GetUserServer() *UserServer {
	UserServerOnce.Do(func() {
		UserServerInstance = &UserServer{}
	})
	return UserServerInstance
}

func (s *UserServer) Hello(ctx context.Context, req *userpb.HelloRequest) (resp *userpb.HelloResponse, err error) {
	resp = &userpb.HelloResponse{
		Msg: "hello",
	}
	return
}

func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (resp *userpb.CreateUserResponse, err error) {
	return nil, err
}

func (s *UserServer) VerifyUser(ctx context.Context, req *userpb.VerifyUserRequest) (resp *userpb.VerifyUserResponse, err error) {
	return nil, err
}

func (s *UserServer) UpdateUserInfo(ctx context.Context, req *userpb.UpdateUserInfoRequest) (resp *userpb.UpdateUserInfoResponse, err error) {
	return nil, err
}

func (s *UserServer) GetUserInfo(ctx context.Context, req *userpb.GetUserInfoRequest) (resp *userpb.GetUserInfoResponse, err error) {
	return nil, err
}

func (s *UserServer) SearchUser(ctx context.Context, req *userpb.SearchUserRequest) (resp *userpb.SearchUserRequest, err error) {
	return nil, err
}
