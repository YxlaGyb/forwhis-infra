package modules

import (
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/admin")
		{
			v1.GET("/reviews", controller.GetPendingReviews)
			v1.POST("/reviews/:review_id/approve", controller.ApproveReview)
			v1.POST("/reviews/:review_id/reject", controller.RejectReview)
		}
	})
}
