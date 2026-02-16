package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/ai")
		v1.Use(middleware.AuthMiddleware())
		{
			v1.GET("/conversations", controller.GetConversations)
			v1.GET("/messages", controller.GetMessages)
			v1.POST("/chat", controller.ChatStream)
			v1.POST("/conversation/del", controller.DeleteConversation)
		}
	})
}
