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
		v1 := r.Group("/api/v1/admin")
		{
			v1.GET("/users", middleware.RequirePerm(constants.PermUserRead)(controller.GetAllUsers))
			v1.POST("/user/:uid/update", middleware.RequirePerm(constants.PermUserEdit)(controller.AdminUpdateUser))
			v1.GET("/reviews", middleware.RequireAuth(controller.GetReviews))
			v1.POST("/reviews/:id/approve", middleware.RequireAuth(controller.ApproveReviewItem))
			v1.POST("/reviews/:id/reject", middleware.RequireAuth(controller.RejectReviewItem))
			v1.GET("/settings", controller.GetOptions)
			v1.POST("/settings", middleware.RequirePerm(constants.PermSettingsEdit)(controller.SetOptions))
		}
	})
}
