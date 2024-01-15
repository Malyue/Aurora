package router

import (
	_debug "Aurora/internal/apps/gateway/router/debug"
	_user "Aurora/internal/apps/gateway/router/user"
	_ws "Aurora/internal/apps/gateway/router/ws"
	"Aurora/internal/apps/gateway/svc"
	"github.com/gin-gonic/gin"
)

func Init(ctx *svc.ServerCtx, r *gin.Engine) {
	_debug.RegisterDebugRoute(ctx, r)
	_user.RegisterUserRouter(ctx, r)
	_ws.RegisterUserRouter(ctx, r)
}
