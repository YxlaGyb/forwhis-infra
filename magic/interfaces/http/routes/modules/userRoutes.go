package modules

import (
	"forestwhisperGo/magic/api/user"
	"forestwhisperGo/magic/controller"
	"forestwhisperGo/magic/interfaces/http/middleware"
	"forestwhisperGo/magic/interfaces/http/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	routes.Register(func(r *gin.Engine) {
		v1 := r.Group("/api/v1/user")
		{
			v1.GET("/users", controller.GetUserByUsername)
			v1.GET("/me", middleware.RequireAuth(user.GetCurrentUserInfo))
			v1.POST("/login", controller.Login)
			v1.POST("/logout", controller.Logout)
			v1.POST("/register", controller.Register)
			v1.POST("/email/code/register", controller.SendEmailCodeRegister)
			v1.POST("/users/del", middleware.RequireAuth(controller.DeleteUserByUid))
			v1.POST("/me/update", middleware.RequireAuth(controller.UpdateProfile))
			v1.POST("/passwd", middleware.RequireAuth(controller.ChangePassword))
		}
	})
}
