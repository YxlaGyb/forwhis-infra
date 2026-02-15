package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/user")
		{
			v1.POST("/login", controller.Login)
		}
	})
}
