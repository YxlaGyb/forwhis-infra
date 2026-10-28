package modules

import (
	"forestwhisperGo/apps/api/internal/handler"
	"forestwhisperGo/apps/api/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/ai")
		v1.Get("/conversations", handler.GetConversations)
		v1.Get("/messages", handler.GetMessages)
		v1.Post("/chat", handler.ChatStream)
		v1.Post("/conversation/del", handler.DeleteConversation)
	})
}
