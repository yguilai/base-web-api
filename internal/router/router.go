package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wam-lab/base-web-api/common/conno"
	v1 "github.com/wam-lab/base-web-api/internal/api/v1"
	"github.com/wam-lab/base-web-api/internal/global"
	"github.com/wam-lab/base-web-api/internal/middleware"
	"time"
)

func Router() *gin.Engine {
	if global.Config.GetString("mode") == conno.PRO {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	r.Use(
		middleware.LoggerWithZap(global.Log, time.RFC3339, true),
		middleware.RecoveryWithZap(global.Log, true),
		middleware.Cors(),
	)

	apiGroup := r.Group("/api/v1")
	v1.InitAuthRouter(apiGroup)
	v1.InitUserRouter(apiGroup)
	return r
}
