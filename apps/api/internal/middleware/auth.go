package middleware

import (
	"forestwhisperGo/apps/api/internal/models"
	rbacSvc "forestwhisperGo/apps/api/internal/service/rbac"
	"forestwhisperGo/apps/api/pkg/constants"
	"forestwhisperGo/apps/api/pkg/db"
	"forestwhisperGo/apps/api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("forestwhisper")
		if token == "" {
			return JSONError(c, "未登录或登录已过期喵", nil)
		}

		payload, err := utils.ValidateToken(token)
		if err != nil {
			return JSONError(c, "Token无效喵", nil)
		}

		if payload.ISS != "lmoadll" || payload.AUD != "lmoadll" {
			return JSONError(c, "Token无效喵", nil)
		}

		c.Locals("uid", payload.UID)
		return c.Next()
	}
}

func PermMiddleware(perm string) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := GetCurrentUser(c)
		if user == nil {
			return JSONError(c, "未登录或登录已过期喵", nil)
		}

		if !rbacSvc.HasPermission(user, perm) {
			return JSONError(c, "没有权限喵", nil)
		}
		return c.Next()
	}
}

func RoleMiddleware(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := GetCurrentUser(c)
		if user == nil {
			return JSONError(c, "未登录或登录已过期喵", nil)
		}

		hasRole := false
		for _, role := range roles {
			if rbacSvc.HasRole(user, role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return JSONError(c, "没有权限喵", nil)
		}
		return c.Next()
	}
}

func RequireAuth(handler fiber.Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := GetCurrentUser(c)
		if user == nil {
			return JSONError(c, "未登录或登录已过期喵", nil)
		}
		c.Locals("currentUser", user)
		return handler(c)
	}
}

func RequirePerm(perm constants.Permission) func(fiber.Handler) fiber.Handler {
	return func(handler fiber.Handler) fiber.Handler {
		return func(c fiber.Ctx) error {
			user := GetCurrentUser(c)
			if user == nil {
				return JSONError(c, "未登录或登录已过期喵", nil)
			}
			if !rbacSvc.HasPermission(user, string(perm)) {
				return JSONError(c, "没有权限喵", nil)
			}
			c.Locals("currentUser", user)
			return handler(c)
		}
	}
}

func RequireRole(roles ...constants.Role) func(fiber.Handler) fiber.Handler {
	return func(handler fiber.Handler) fiber.Handler {
		return func(c fiber.Ctx) error {
			user := GetCurrentUser(c)
			if user == nil {
				return JSONError(c, "未登录或登录已过期喵", nil)
			}
			hasRole := false
			for _, role := range roles {
				if rbacSvc.HasRole(user, string(role)) {
					hasRole = true
					break
				}
			}
			if !hasRole {
				return JSONError(c, "没有权限喵", nil)
			}
			c.Locals("currentUser", user)
			return handler(c)
		}
	}
}

func GetCurrentUser(c fiber.Ctx) *models.User {
	if uid, ok := c.Locals("uid").(int); ok {
		var user models.User
		if db.GetDB().First(&user, uid).Error == nil {
			return &user
		}
	}

	token := c.Cookies("forestwhisper")
	if token == "" {
		return nil
	}

	payload, err := utils.ValidateToken(token)
	if err != nil || payload.ISS != "lmoadll" || payload.AUD != "lmoadll" {
		return nil
	}

	var user models.User
	if db.GetDB().First(&user, payload.UID).Error != nil {
		return nil
	}

	c.Locals("uid", user.Uid)
	return &user
}

func GetCurrentUserFromContext(c fiber.Ctx) *models.User {
	if user, ok := c.Locals("currentUser").(*models.User); ok {
		return user
	}
	return GetCurrentUser(c)
}

func GetUserWithRoles(c fiber.Ctx) *models.User {
	user := GetCurrentUser(c)
	if user == nil {
		return nil
	}
	return user
}
