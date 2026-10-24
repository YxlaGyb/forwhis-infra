package modules

import (
	"forestwhisperGo/apps/internal/controller"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/content")
		v1.Get("", controller.GetContents)
		v1.Get("/my", middleware.RequireAuth(controller.GetMyContents))
		v1.Get("/:cid", controller.GetContent)
		v1.Post("", middleware.RequireAuth(controller.CreateContent))
		v1.Post("/:cid", middleware.RequireAuth(controller.UpdateContent))
		v1.Post("/:cid/del", middleware.RequireAuth(controller.DeleteContent))
	})
}
