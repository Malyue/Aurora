package _grpc

import (
	userpb "Aurora/api/proto-go/user"
	_const "Aurora/internal/pkg/const"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var UserServiceClient userpb.UserServiceClient

func InitUserClient() (userpb.UserServiceClient, error) {
	conn, err := initClient(_const.UserServiceName)
	if err != nil {
		return nil, err
	}
	UserServiceClient = userpb.NewUserServiceClient(conn)
	return UserServiceClient, nil
}

func initClient(name string) (*grpc.ClientConn, error) {
	target := fmt.Sprintf("etcd:///%s", name)

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return conn, err
}
