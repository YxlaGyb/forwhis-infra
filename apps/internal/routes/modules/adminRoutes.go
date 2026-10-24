package modules

import (
	"forestwhisperGo/apps/internal/controller"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"
	"forestwhisperGo/apps/pkg/constants"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/admin")
		v1.Get("/users", middleware.RequirePerm(constants.PermUserRead)(controller.GetAllUsers))
		v1.Post("/user/create", middleware.RequirePerm(constants.PermUserEdit)(controller.AdminCreateUser))
		v1.Post("/user/:uid/update", middleware.RequirePerm(constants.PermUserEdit)(controller.AdminUpdateUser))
		v1.Get("/reviews", middleware.RequireAuth(controller.GetReviews))
		v1.Post("/reviews/:id/approve", middleware.RequireAuth(controller.ApproveReviewItem))
		v1.Post("/reviews/:id/reject", middleware.RequireAuth(controller.RejectReviewItem))
		v1.Get("/settings", controller.GetOptions)
		v1.Post("/settings", middleware.RequirePerm(constants.PermSettingsEdit)(controller.SetOptions))
	})
}
