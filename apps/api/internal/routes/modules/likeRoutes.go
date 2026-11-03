package modules

import (
	"forwhis-infra/apps/api/internal/handler"
	"forwhis-infra/apps/api/internal/middleware"
	"forwhis-infra/apps/api/internal/routes"

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
