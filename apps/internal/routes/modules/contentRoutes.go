package modules

import (
	"forestwhisperGo/apps/internal/handler"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/content")
		v1.Get("", handler.GetContents)
		v1.Get("/my", middleware.RequireAuth(handler.GetMyContents))
		v1.Get("/:cid", handler.GetContent)
		v1.Post("", middleware.RequireAuth(handler.CreateContent))
		v1.Post("/:cid", middleware.RequireAuth(handler.UpdateContent))
		v1.Post("/:cid/del", middleware.RequireAuth(handler.DeleteContent))
	})
}
