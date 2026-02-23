package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/comment")
		{
			v1.GET("", controller.GetComments)
			v1.GET("/my", middleware.RequireAuth(controller.GetMyComments))
			v1.GET("/content/:cid", controller.GetContentComments)
			v1.GET("/:coid", controller.GetComment)
			v1.POST("", controller.CreateComment)
			v1.POST("/:coid", middleware.RequireAuth(controller.UpdateComment))
			v1.POST("/:coid/approve", middleware.RequireAuth(controller.ApproveComment))
			v1.POST("/:coid/del", middleware.RequireAuth(controller.DeleteComment))
		}
	})
}
