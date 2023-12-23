package user

import (
	userpb "Aurora/api/proto-go/user"
	userdefine "Aurora/internal/apps/geteway/define/user"
	"Aurora/internal/apps/geteway/svc"
	"Aurora/internal/pkg/errorx"
	"Aurora/internal/pkg/resp"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req userdefine.RegisterRequest
		if err := ctx.ShouldBind(&req); err != nil {
			svcCtx.Logger.Errorf("Bind request %s", err)
			resp.ResponseError(ctx, errorx.CodeInvalidParam)
			return
		}

		// TODO valid param
		response, err := svcCtx.UserServer.CreateUser(ctx, &userpb.CreateUserRequest{
			Account:  req.Account,
			Password: req.Password,
		})
		if err != nil {
			svcCtx.Logger.Errorf("User Rpc request error : %s", err)
			resp.ResponseError(ctx, errorx.CodeServerBusy)
			return
		}

		// if it is no success, it means it is repeated account
		if !response.Success {
			resp.ResponseError(ctx, errorx.CodeRepeatedAccount)
			return
		}

		resp.ResponseSuccess(ctx, &userdefine.RegisterResponse{})
	}
}
