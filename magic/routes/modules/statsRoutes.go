package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/admin/stats")
		{
			v1.GET("/dashboard", controller.GetDashboard)
			v1.GET("/users", controller.GetUserStats)
			v1.GET("/contents", controller.GetContentStats)
			v1.GET("/comments", controller.GetCommentStats)
			v1.GET("/trend", controller.GetTrendStats)
		}
	})
}
