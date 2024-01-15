package ws

import (
	_ws "Aurora/internal/apps/gateway/handler/ws"
	"Aurora/internal/apps/gateway/svc"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(ctx *svc.ServerCtx, router *gin.Engine) {
	ws := router.Group("/ws")
	{
		ws.GET("/", _ws.WebSocketProxyHandler(ctx))
	}
}
