package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/interfaces/http/middleware"
	"forestwhisperGo/magic/interfaces/http/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/meta")
		{
			v1.GET("", controller.GetMetas)
			v1.GET("/:mid", controller.GetMeta)
			v1.POST("", middleware.RequireAuth(controller.CreateMeta))
			v1.POST("/:mid", middleware.RequireAuth(controller.UpdateMeta))
			v1.POST("/:mid/del", middleware.RequireAuth(controller.DeleteMeta))
		}
	})
}
