package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/interfaces/http/middleware"
	"forestwhisperGo/magic/interfaces/http/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/content")
		{
			v1.GET("", controller.GetContents)
			v1.GET("/my", middleware.RequireAuth(controller.GetMyContents))
			v1.GET("/:cid", controller.GetContent)
			v1.POST("", middleware.RequireAuth(controller.CreateContent))
			v1.POST("/:cid", middleware.RequireAuth(controller.UpdateContent))
			v1.POST("/:cid/del", middleware.RequireAuth(controller.DeleteContent))
		}
	})
}
