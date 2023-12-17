package user

import (
	userpb "Aurora/api/proto-go/user"
	"Aurora/internal/apps/geteway/svc"
	"github.com/gin-gonic/gin"
)

func HelloHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		msg, _ := svcCtx.UserServer.Hello(ctx, &userpb.HelloRequest{})
		ctx.JSON(200, msg)
	}
}
