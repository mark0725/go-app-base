package auth

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(group string, r *gin.RouterGroup) {
	switch group {
	case "sapi":
	case "api":
		r.POST("/login", g_AuthApi.Login)
		r.GET("/logout", g_AuthApi.Loginout)
	case "*":

	default:

	}
}
