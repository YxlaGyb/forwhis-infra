package modules

import (
	"forestwhisperGo/apps/internal/controller"
	"forestwhisperGo/apps/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/admin/stats")
		v1.Get("/dashboard", controller.GetDashboard)
		v1.Get("/users", controller.GetUserStats)
		v1.Get("/contents", controller.GetContentStats)
		v1.Get("/comments", controller.GetCommentStats)
		v1.Get("/trend", controller.GetTrendStats)
	})
}
