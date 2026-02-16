package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/content")
		{
			v1.GET("/", controller.GetContents)
			v1.GET("/my", controller.GetMyContents)
			v1.GET("/:cid", controller.GetContent)
			v1.POST("/", controller.CreateContent)
			v1.POST("/:cid", controller.UpdateContent)
			v1.POST("/:cid/del", controller.DeleteContent)
		}
	})
}
