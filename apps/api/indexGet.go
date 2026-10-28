package user

import (
	"net/http"

	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/service"
	rbacSvc "forestwhisperGo/apps/api/internal/service/rbac"

	"github.com/gofiber/fiber/v3"
)

func GetCurrentUserInfo(c fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return c.Status(http.StatusOK).JSON(middleware.Error("未登录或登录已过期喵喵", nil))
	}

	profile := service.GetUserProfile(user.Uid)
	roles := rbacSvc.GetUserRoles(user.Uid)
	canAccessAdmin := rbacSvc.HasAnyAdminRole(user)

	data := fiber.Map{
		"uid":            user.Uid,
		"name":           user.Name,
		"email":          user.Mail,
		"lastip":         user.LastIp,
		"createdAt":      user.CreatedAt,
		"lastLogin":      user.LastLogin,
		"isActive":       user.IsActive,
		"roles":          roles,
		"canAccessAdmin": canAccessAdmin,
	}
	for k, v := range profile {
		data[k] = v
	}
	return c.Status(http.StatusOK).JSON(middleware.Success("ok", data))
}
