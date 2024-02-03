package user

import (
	userhandler "Aurora/internal/apps/gateway/handler/user"
	"Aurora/internal/apps/gateway/svc"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(ctx *svc.ServerCtx, router *gin.Engine) {
	user := router.Group("/api/v1/user")
	user.POST("/register", userhandler.RegisterHandler(ctx))
	user.POST("/login", userhandler.LoginHandler(ctx))
}
