package fix

import (
	"forestwhisperGo/magic/interfaces/http/middleware"
	"forestwhisperGo/pkg/constants"
	"forestwhisperGo/pkg/utils"

	"github.com/gin-gonic/gin"
)

func SetupFixRoutes(r *gin.Engine) {
	if !isFixToolEnabled() {
		return
	}

	fixApi := r.Group("/api/fix")
	fixApi.Use(middleware.AuthMiddleware())
	fixApi.Use(middleware.RoleMiddleware(string(constants.RoleSuperadmin), string(constants.RoleAdmin)))
	{
		fixApi.GET("/status", FixStatusHandler)
		fixApi.GET("/tables", CheckTablesHandler)
		fixApi.GET("/table/info", GetTableInfoHandler)
		fixApi.POST("/tables/fix", FixTablesHandler)
		fixApi.POST("/table/fix", FixSingleTableHandler)
		fixApi.POST("/indexes", CreateIndexesHandler)
	}
}

func isFixToolEnabled() bool {
	enabled := utils.GetConfig("fix", "enabled")
	if enabled == nil {
		return false
	}
	if b, ok := enabled.(bool); ok {
		return b
	}
	return false
}
