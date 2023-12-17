package router

import (
	_debug "Aurora/internal/apps/geteway/router/debug"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	_debug.RegisterDebugRoute(r)
}
