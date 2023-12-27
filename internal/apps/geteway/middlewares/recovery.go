package middleware

import (
	"Aurora/internal/pkg/errorx"
	"Aurora/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Info(string(debug.Stack()))
				resp.ResponseError(c, errorx.CodeServerBusy)
				return
			}
		}()
		c.Next()
	}
}
