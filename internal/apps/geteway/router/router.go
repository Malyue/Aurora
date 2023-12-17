package router

import (
	"Aurora/internal/apps/geteway/handler/user"
	middleware "Aurora/internal/apps/geteway/middlewares"
	"Aurora/internal/apps/geteway/svc"
	"github.com/gin-gonic/gin"
)

func InitRouter(ctx *svc.ServerCtx) *gin.Engine {
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

	r.GET("/hello", user.HelloHandler(ctx))

	Init(ctx, r)

	return r
}
