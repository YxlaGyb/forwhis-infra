package modules

import (
	"forestwhisperGo/magic/api/user"
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/user")
		{
			v1.GET("/users", controller.GetUserByUsername)
			v1.GET("/me", user.GetCurrentUserInfo)
			v1.POST("/login", controller.Login)
			v1.POST("/register", controller.Register)
			v1.POST("/email/code/register", controller.SendEmailCodeRegister)
			v1.POST("/users/del", controller.DeleteUserByUid)
			v1.POST("/me/update", controller.UpdateProfile)
			v1.POST("/passwd", controller.ChangePassword)
		}
	})
}
