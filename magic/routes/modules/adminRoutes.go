package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/admin")
		{
			v1.GET("/reviews", middleware.RequireAuth(controller.GetReviews))
			v1.POST("/reviews/:id/approve", middleware.RequireAuth(controller.ApproveReviewItem))
			v1.POST("/reviews/:id/reject", middleware.RequireAuth(controller.RejectReviewItem))
		}
	})
}
