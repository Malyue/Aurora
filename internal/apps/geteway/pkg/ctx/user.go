package ctx

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_const "Aurora/internal/pkg/const"
	"Aurora/internal/pkg/errorx"
)

func GetUserIDInCtx(ctx *gin.Context) (string, error) {
	userId, ok := ctx.Get(_const.UserIDCtx)
	if !ok {
		logrus.Error(_const.GateWay, "Get UserID in ctx error")
		return "", errorx.CodeErrAuth
	}
	return userId.(string), nil
}
