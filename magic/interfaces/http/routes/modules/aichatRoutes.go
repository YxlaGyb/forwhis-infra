package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/interfaces/http/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/ai")
		{
			v1.GET("/conversations", controller.GetConversations)
			v1.GET("/messages", controller.GetMessages)
			v1.POST("/chat", controller.ChatStream)
			v1.POST("/conversation/del", controller.DeleteConversation)
		}
	})
}
