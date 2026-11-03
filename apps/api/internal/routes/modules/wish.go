package modules

import (
	"forwhis-infra/apps/api/internal/handler"
	"forwhis-infra/apps/api/internal/middleware"
	"forwhis-infra/apps/api/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/wish")
		v1.Get("", handler.GetWishes)
		v1.Get("/my", middleware.RequireAuth(handler.GetMyWishes))
		v1.Post("", handler.CreateWish)
		v1.Post("/:id/del", middleware.RequireAuth(handler.DeleteWish))
	})
}
