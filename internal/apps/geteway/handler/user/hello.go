package user

import (
	userpb "Aurora/api/proto-go/user"
	_grpc "Aurora/internal/pkg/grpc"
	"github.com/gin-gonic/gin"
)

func HelloHandler(ctx *gin.Context) {
	msg, _ := _grpc.UserServiceClient.Hello(ctx, &userpb.HelloRequest{})
	ctx.JSON(200, msg)
}
