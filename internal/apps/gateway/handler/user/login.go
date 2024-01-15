package user

import (
	userpb "Aurora/api/proto-go/user"
	userdefine "Aurora/internal/apps/gateway/define/user"
	"Aurora/internal/apps/gateway/svc"
	"Aurora/internal/pkg/errorx"
	"Aurora/internal/pkg/jwt"
	"Aurora/internal/pkg/resp"
	"github.com/gin-gonic/gin"
)

func LoginHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req userdefine.LoginRequest
		if err := ctx.ShouldBind(&req); err != nil {
			svcCtx.Logger.Errorf("Bind request %s", err)
			resp.ResponseError(ctx, errorx.CodeInvalidParam)
			return
		}

		// TODO valid param
		response, err := svcCtx.UserServer.VerifyUser(ctx, &userpb.VerifyUserRequest{
			Account:  req.Account,
			Password: req.Password,
		})
		if err != nil {
			svcCtx.Logger.Errorf("User Rpc request error : %s", err)
			resp.ResponseError(ctx, errorx.CodeServerBusy)
			return
		}
		if !response.Verify {
			resp.ResponseError(ctx, errorx.CodeErrPassword)
			return
		}

		// get accessToken and refreshToken
		accessToken, refreshToken, err := jwt.GenerateToken(response.User.Id)
		if err != nil {
			svcCtx.Logger.Errorf("Generate token  error : %s in id : %s", err, response.User.Id)
			resp.ResponseError(ctx, errorx.CodeServerBusy)
			return
		}

		resp.ResponseSuccess(ctx, &userdefine.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
