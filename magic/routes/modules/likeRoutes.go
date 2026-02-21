package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/like")
		{
			v1.POST("/toggle", middleware.RequireAuth(controller.ToggleLike))
			v1.GET("/count", controller.GetLikeCount)
			v1.GET("/status", controller.GetLikeStatus)
			v1.GET("/user", middleware.RequireAuth(controller.GetUserLikes))
		}
	})
}
