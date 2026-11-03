package modules

import (
	"forwhis-infra/apps/api/internal/handler"
	"forwhis-infra/apps/api/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/feed")
		v1.Get("", handler.GetFeed)
	})
}
