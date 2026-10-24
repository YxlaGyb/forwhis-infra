package modules

import (
	"forestwhisperGo/apps/internal/controller"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/like")
		v1.Post("/toggle", middleware.RequireAuth(controller.ToggleLike))
		v1.Get("/count", controller.GetLikeCount)
		v1.Get("/status", controller.GetLikeStatus)
		v1.Get("/user", middleware.RequireAuth(controller.GetUserLikes))
		v1.Get("/my/count", middleware.RequireAuth(controller.GetMyLikesCount))
	})
}
