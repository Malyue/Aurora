package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"Aurora/internal/pkg/errorx"
)

type ResponseCode struct {
	Code errorx.ResCode `json:"code"`
	Msg  interface{}    `json:"msg"`
	Data interface{}    `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code errorx.ResCode) {
	c.JSON(http.StatusOK, &ResponseCode{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseCode{
		Code: errorx.CodeSuccess,
		Data: data,
	})
}

func ResponseWithMsg(c *gin.Context, code errorx.ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseCode{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
