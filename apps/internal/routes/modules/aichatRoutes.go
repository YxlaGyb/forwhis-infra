package modules

import (
	"forestwhisperGo/apps/internal/controller"
	"forestwhisperGo/apps/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/ai")
		v1.Get("/conversations", controller.GetConversations)
		v1.Get("/messages", controller.GetMessages)
		v1.Post("/chat", controller.ChatStream)
		v1.Post("/conversation/del", controller.DeleteConversation)
	})
}
