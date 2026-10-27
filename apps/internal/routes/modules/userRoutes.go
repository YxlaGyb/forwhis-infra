package modules

import (
	user "forestwhisperGo/apps"
	"forestwhisperGo/apps/internal/handler"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func init() {
	routes.Register(func(r *fiber.App) {
		v1 := r.Group("/api/v1/user")
		v1.Get("/users", handler.GetUserByUsername)
		v1.Get("/me", middleware.RequireAuth(user.GetCurrentUserInfo))
		v1.Post("/login", handler.Login)
		v1.Post("/logout", handler.Logout)
		v1.Post("/register", handler.Register)
		v1.Post("/email/code/register", handler.SendEmailCodeRegister)
		v1.Post("/users/del", middleware.RequireAuth(handler.DeleteUserByUid))
		v1.Post("/me/update", middleware.RequireAuth(handler.UpdateProfile))
		v1.Post("/passwd", middleware.RequireAuth(handler.ChangePassword))
	})
}
