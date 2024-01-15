package user

import (
	"github.com/gin-gonic/gin"

	userdefine "Aurora/internal/apps/gateway/define/user"
	_ctx "Aurora/internal/apps/gateway/pkg/ctx"
	"Aurora/internal/apps/gateway/svc"
	"Aurora/internal/pkg/errorx"
	"Aurora/internal/pkg/resp"
)

func GetUserInfo(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get user id in ctx
		userId, err := _ctx.GetUserIDInCtx(ctx)
		if err != nil {
			resp.ResponseError(ctx, errorx.CodeErrAuth)
			return
		}

		// TODO construct Req
		_ = userdefine.GetUserInfoRequest{
			ID: userId,
		}

	}
}
