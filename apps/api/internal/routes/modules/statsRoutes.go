package modules

import (
	"forestwhisperGo/apps/api/internal/handler"
	"forestwhisperGo/apps/api/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/admin/stats")
		v1.Get("/dashboard", handler.GetDashboard)
		v1.Get("/users", handler.GetUserStats)
		v1.Get("/contents", handler.GetContentStats)
		v1.Get("/comments", handler.GetCommentStats)
		v1.Get("/trend", handler.GetTrendStats)
	})
}
