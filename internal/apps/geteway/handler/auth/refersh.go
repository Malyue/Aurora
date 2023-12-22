package auth

import (
	"Aurora/internal/apps/geteway/define/auth"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"Aurora/internal/apps/geteway/svc"
	_const "Aurora/internal/pkg/const"
	"Aurora/internal/pkg/errorx"
	"Aurora/internal/pkg/jwt"
	"Aurora/internal/pkg/resp"
)

// RefreshToken when accessToken is expired, valid the refreshToken
// if refreshToken is not expired, refresh the accessToken and refreshToken
func RefreshToken(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get refreshToken
		token := ctx.GetHeader(_const.RefreshToken)

		// parse token
		accessToken, refreshToken, err := jwt.ParseRefreshToken(token)
		if err != nil {
			resp.ResponseError(ctx, errorx.CodeErrAuth)
			logrus.Error(_const.GateWay, err)
			return
		}

		resp.ResponseSuccess(ctx, &auth.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
