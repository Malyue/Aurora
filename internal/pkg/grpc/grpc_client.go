package _grpc

import (
	userpb "Aurora/api/proto-go/user"
	_const "Aurora/internal/pkg/const"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var UserServiceClient userpb.UserServiceClient

func InitUserClient() userpb.UserServiceClient {
	conn := initClient(_const.UserServiceName)
	UserServiceClient = userpb.NewUserServiceClient(conn)
	return UserServiceClient
}

func initClient(name string) *grpc.ClientConn {
	target := fmt.Sprintf("etcd:///%s", name)

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	return conn
}
