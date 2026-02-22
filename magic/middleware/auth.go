package middleware

import (
	"forestwhisperGo/magic/constants"
	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/rbac"
	rbacSvc "forestwhisperGo/magic/service/rbac"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("forestwhisper")
		if err != nil || token == "" {
			JSONError(c, "未登录或登录已过期喵", nil)
			c.Abort()
			return
		}

		payload, err := utils.ValidateToken(token)
		if err != nil {
			JSONError(c, "Token无效喵", nil)
			c.Abort()
			return
		}

		if payload.ISS != "lmoadll" || payload.AUD != "lmoadll" {
			JSONError(c, "Token无效喵", nil)
			c.Abort()
			return
		}

		c.Set("uid", payload.UID)
		c.Next()
	}
}

func PermMiddleware(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			JSONError(c, "未登录或登录已过期喵", nil)
			c.Abort()
			return
		}

		if !rbacSvc.HasPermission(user, perm) {
			JSONError(c, "没有权限喵", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			JSONError(c, "未登录或登录已过期喵", nil)
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if rbacSvc.HasPermission(user, role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			JSONError(c, "没有权限喵", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAuth(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			JSONError(c, "未登录或登录已过期喵", nil)
			c.Abort()
			return
		}
		c.Set("currentUser", user)
		handler(c)
	}
}

func RequirePerm(perm constants.Permission) func(gin.HandlerFunc) gin.HandlerFunc {
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			user := GetCurrentUser(c)
			if user == nil {
				JSONError(c, "未登录或登录已过期喵", nil)
				c.Abort()
				return
			}
			if !rbacSvc.HasPermission(user, string(perm)) {
				JSONError(c, "没有权限喵", nil)
				c.Abort()
				return
			}
			c.Set("currentUser", user)
			handler(c)
		}
	}
}

func RequireRole(roles ...constants.Role) func(gin.HandlerFunc) gin.HandlerFunc {
	return func(handler gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			user := GetCurrentUser(c)
			if user == nil {
				JSONError(c, "未登录或登录已过期喵", nil)
				c.Abort()
				return
			}
			hasRole := false
			for _, role := range roles {
				if rbacSvc.HasPermission(user, string(role)) {
					hasRole = true
					break
				}
			}
			if !hasRole {
				JSONError(c, "没有权限喵", nil)
				c.Abort()
				return
			}
			c.Set("currentUser", user)
			handler(c)
		}
	}
}

func GetCurrentUser(c *gin.Context) *models.User {
	if uid, exists := c.Get("uid"); exists {
		var user models.User
		if db.GetDB().First(&user, uid).Error == nil {
			return &user
		}
	}

	token, err := c.Cookie("forestwhisper")
	if err != nil || token == "" {
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

	c.Set("uid", user.Uid)
	return &user
}

func GetCurrentUserFromContext(c *gin.Context) *models.User {
	if user, exists := c.Get("currentUser"); exists {
		return user.(*models.User)
	}
	return GetCurrentUser(c)
}

func GetUserWithRoles(c *gin.Context) *models.User {
	user := GetCurrentUser(c)
	if user == nil {
		return nil
	}
	var roles []rbac.UserRole
	db.GetDB().Where("user_id = ?", user.Uid).Preload("Role").Find(&roles)
	return user
}
