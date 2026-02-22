package modules

import (
	"forestwhisperGo/magic/constants"
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/options")
		{
			v1.GET("", controller.GetOptions)
			v1.GET("/:key", controller.GetOption)
			v1.POST("", middleware.RequirePerm(constants.PermSettingsEdit)(controller.SetOptions))
			v1.POST("/set", middleware.RequirePerm(constants.PermSettingsEdit)(controller.SetOption))
			v1.POST("/delete/:key", middleware.RequirePerm(constants.PermSettingsEdit)(controller.DeleteOption))
		}
	})
}
