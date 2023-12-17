package router

import (
	_debug "Aurora/internal/apps/geteway/router/debug"
	_user "Aurora/internal/apps/geteway/router/user"
	"Aurora/internal/apps/geteway/svc"
	"github.com/gin-gonic/gin"
)

func Init(ctx *svc.ServerCtx, r *gin.Engine) {
	_debug.RegisterDebugRoute(ctx, r)
	_user.RegisterRouter(ctx, r)
}
