package modules

import (
	"forestwhisperGo/apps/internal/controller"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/meta")
		v1.Get("", controller.GetMetas)
		v1.Get("/:mid", controller.GetMeta)
		v1.Post("", middleware.RequireAuth(controller.CreateMeta))
		v1.Post("/:mid", middleware.RequireAuth(controller.UpdateMeta))
		v1.Post("/:mid/del", middleware.RequireAuth(controller.DeleteMeta))
	})
}
