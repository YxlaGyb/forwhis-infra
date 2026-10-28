package modules

import (
	"forestwhisperGo/apps/api/internal/handler"
	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/routes"
	"forestwhisperGo/apps/api/pkg/constants"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/admin")
		v1.Get("/users", middleware.RequirePerm(constants.PermUserRead)(handler.GetAllUsers))
		v1.Post("/user/create", middleware.RequirePerm(constants.PermUserEdit)(handler.AdminCreateUser))
		v1.Post("/user/:uid/update", middleware.RequirePerm(constants.PermUserEdit)(handler.AdminUpdateUser))
		v1.Get("/reviews", middleware.RequirePerm(constants.PermCommentModerate)(handler.GetReviews))
		v1.Post("/reviews/:id/approve", middleware.RequireAuth(handler.ApproveReviewItem))
		v1.Post("/reviews/:id/reject", middleware.RequireAuth(handler.RejectReviewItem))
		v1.Get("/settings", handler.GetOptions)
		v1.Post("/settings", middleware.RequirePerm(constants.PermSettingsEdit)(handler.SetOptions))
	})
}
