package router

import (
	"Aurora/internal/apps/gateway/handler/auth"
	"Aurora/internal/apps/gateway/handler/user"
	middleware "Aurora/internal/apps/gateway/middlewares"
	"Aurora/internal/apps/gateway/svc"
	"github.com/gin-gonic/gin"
)

func InitRouter(svcCtx *svc.ServerCtx) *gin.Engine {
	r := gin.New()
	r.Use(middleware.CorsMiddleware())

	// TODO swagger
	//r.GET("/swagger/*any",gx.WrapHandler(swa))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, map[string]any{"code": 200, "msg": "hello world"})
	})

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, map[string]any{"status": "ok"})
	})

	r.GET("/hello", user.HelloHandler(svcCtx))

	// register refresh token handler
	r.GET("/refreshToken", auth.RefreshToken(svcCtx))

	Init(svcCtx, r)

	return r
}
