package modules

import (
	"forestwhisperGo/apps/api/internal/handler"
	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/comment")
		v1.Get("", handler.GetComments)
		v1.Get("/my", middleware.RequireAuth(handler.GetMyComments))
		v1.Get("/content/:cid", handler.GetContentComments)
		v1.Get("/:coid", handler.GetComment)
		v1.Post("", handler.CreateComment)
		v1.Post("/:coid", middleware.RequireAuth(handler.UpdateComment))
		v1.Post("/:coid/approve", middleware.RequireAuth(handler.ApproveComment))
		v1.Post("/:coid/del", middleware.RequireAuth(handler.DeleteComment))
	})
}
