package fix

import (
	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/pkg/constants"
	"forestwhisperGo/apps/api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

func SetupFixRoutes(r *fiber.App) {
	if !isFixToolEnabled() {
		return
	}

	fixApi := r.Group("/api/fix")
	fixApi.Use(middleware.AuthMiddleware())
	fixApi.Use(middleware.PermMiddleware(string(constants.PermSystemConfig)))
	{
		fixApi.Get("/status", FixStatusHandler)
		fixApi.Get("/tables", CheckTablesHandler)
		fixApi.Get("/table/info", GetTableInfoHandler)
		fixApi.Post("/tables/fix", FixTablesHandler)
		fixApi.Post("/table/fix", FixSingleTableHandler)
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
