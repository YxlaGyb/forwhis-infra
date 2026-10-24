package modules

import (
	"forestwhisperGo/apps/internal/controller"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/comment")
		v1.Get("", controller.GetComments)
		v1.Get("/my", middleware.RequireAuth(controller.GetMyComments))
		v1.Get("/content/:cid", controller.GetContentComments)
		v1.Get("/:coid", controller.GetComment)
		v1.Post("", controller.CreateComment)
		v1.Post("/:coid", middleware.RequireAuth(controller.UpdateComment))
		v1.Post("/:coid/approve", middleware.RequireAuth(controller.ApproveComment))
		v1.Post("/:coid/del", middleware.RequireAuth(controller.DeleteComment))
	})
}
