package middleware

import (
	_const "Aurora/internal/pkg/const"
	"Aurora/internal/pkg/errorx"
	"Aurora/internal/pkg/jwt"
	"Aurora/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"strings"
)

// AuthMiddleware auth
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		headerList := strings.Split(header, " ")
		if len(headerList) != 2 {
			resp.ResponseError(c, errorx.CodeErrAuth)
			c.Abort()
			return
		}

		tokenType := headerList[0]
		token := headerList[1]
		if tokenType != "Bearer " {
			resp.ResponseError(c, errorx.CodeErrAuth)
			c.Abort()
			return
		}

		claims, expire, err := jwt.ParseTokenAndValidExpire(token)
		if err != nil {
			resp.ResponseError(c, errorx.CodeErrAuth)
			c.Abort()
			return
		}

		if expire {
			resp.ResponseError(c, errorx.CodeTokenExpire)
			c.Abort()
			return
		}
		// set id in ctx
		c.Set(_const.UserIDCtx, claims.Id)
		c.Next()
	}
}
