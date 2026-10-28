package modules

import (
	"forestwhisperGo/apps/api/internal/handler"
	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/routes"
	"forestwhisperGo/apps/api/pkg/constants"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/options")
		v1.Get("", handler.GetOptions)
		v1.Get("/:key", handler.GetOption)
		v1.Post("", middleware.RequirePerm(constants.PermSettingsEdit)(handler.SetOptions))
		v1.Post("/set", middleware.RequirePerm(constants.PermSettingsEdit)(handler.SetOption))
		v1.Post("/delete/:key", middleware.RequirePerm(constants.PermSettingsEdit)(handler.DeleteOption))
	})
}
