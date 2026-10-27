package modules

import (
	"forestwhisperGo/apps/internal/handler"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/like")
		v1.Post("/toggle", middleware.RequireAuth(handler.ToggleLike))
		v1.Get("/count", handler.GetLikeCount)
		v1.Get("/status", handler.GetLikeStatus)
		v1.Get("/user", middleware.RequireAuth(handler.GetUserLikes))
		v1.Get("/my/count", middleware.RequireAuth(handler.GetMyLikesCount))
	})
}
