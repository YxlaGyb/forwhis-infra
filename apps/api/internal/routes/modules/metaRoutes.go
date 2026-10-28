package modules

import (
	"forestwhisperGo/apps/api/internal/handler"
	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/meta")
		v1.Get("", handler.GetMetas)
		v1.Get("/:mid", handler.GetMeta)
		v1.Post("", middleware.RequireAuth(handler.CreateMeta))
		v1.Post("/:mid", middleware.RequireAuth(handler.UpdateMeta))
		v1.Post("/:mid/del", middleware.RequireAuth(handler.DeleteMeta))
	})
}
