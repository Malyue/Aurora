package middleware

import (
	userpb "Aurora/api/proto-go/user"
	"Aurora/internal/apps/gateway/svc"
	_const "Aurora/internal/pkg/const"
	"Aurora/internal/pkg/errorx"
	"Aurora/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"strings"
)

// AuthMiddleware auth
func AuthMiddleware(svcCtx *svc.ServerCtx) gin.HandlerFunc {
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

		verifyTokenResp, err := svcCtx.UserServer.VerifyToken(c, &userpb.VerifyTokenRequest{
			Token: token,
		})
		if err != nil {
			resp.ResponseError(c, errorx.CodeErrAuth)
			c.Abort()
			return
		}

		if verifyTokenResp.Expire {
			resp.ResponseError(c, errorx.CodeTokenExpire)
			c.Abort()
			return
		}
		// set id in ctx
		c.Set(_const.UserIDCtx, verifyTokenResp.Id)
		c.Next()
	}
}
