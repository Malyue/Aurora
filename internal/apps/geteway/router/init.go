package router

import (
	_debug "Aurora/internal/apps/geteway/router/debug"
	_user "Aurora/internal/apps/geteway/router/user"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	_debug.RegisterDebugRoute(r)
	_user.RegisterRouter(r)
}
