package modules

import (
	"forestwhisperGo/apps"
	"forestwhisperGo/apps/internal/controller"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/user")
		v1.Get("/users", controller.GetUserByUsername)
		v1.Get("/me", middleware.RequireAuth(user.GetCurrentUserInfo))
		v1.Post("/login", controller.Login)
		v1.Post("/logout", controller.Logout)
		v1.Post("/register", controller.Register)
		v1.Post("/email/code/register", controller.SendEmailCodeRegister)
		v1.Post("/users/del", middleware.RequireAuth(controller.DeleteUserByUid))
		v1.Post("/me/update", middleware.RequireAuth(controller.UpdateProfile))
		v1.Post("/passwd", middleware.RequireAuth(controller.ChangePassword))
	})
}
