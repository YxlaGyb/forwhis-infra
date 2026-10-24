package modules

import (
	"forestwhisperGo/apps/internal/controller"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"
	"forestwhisperGo/apps/pkg/constants"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/options")
		v1.Get("", controller.GetOptions)
		v1.Get("/:key", controller.GetOption)
		v1.Post("", middleware.RequirePerm(constants.PermSettingsEdit)(controller.SetOptions))
		v1.Post("/set", middleware.RequirePerm(constants.PermSettingsEdit)(controller.SetOption))
		v1.Post("/delete/:key", middleware.RequirePerm(constants.PermSettingsEdit)(controller.DeleteOption))
	})
}
