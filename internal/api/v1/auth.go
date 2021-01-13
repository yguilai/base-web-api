package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wam-lab/base-web-api/common/errno"
	"github.com/wam-lab/base-web-api/internal/global"
	"github.com/wam-lab/base-web-api/internal/global/response"
	"github.com/wam-lab/base-web-api/internal/middleware"
)

func InitAuthRouter(g *gin.RouterGroup) {
	auth := g.Group("/auth")
	{
		auth.POST("/login", Login)
		auth.POST("/register", Register)
	}

	auth.Use(middleware.JwtAuth())
	{
		auth.POST("/refresh", RefreshToken)
	}
}

func Login(c *gin.Context) {
	global.Log.Info("Login APi")
	response.Json(c, errno.OK)
}

func Register(c *gin.Context) {
	global.Log.Info("Register APi")
	response.Json(c, errno.OK)
}

func RefreshToken(c *gin.Context) {
	response.Json(c, errno.OK.WithData(map[string]interface{}{
		"newToken": "xxx.xx.xx",
	}))
}