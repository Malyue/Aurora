package ws

import (
	"Aurora/internal/apps/gateway/svc"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func WebSocketProxyHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println(111)
		svcCtx.Logger.Info("conn")
		target, _ := url.Parse("ws://localhost:8082")
		proxy := httputil.NewSingleHostReverseProxy(target)
		director := proxy.Director
		proxy.Director = func(req *http.Request) {
			director(req)
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
			req.Host = req.URL.Host
		}
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
