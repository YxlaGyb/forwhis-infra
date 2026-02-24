package user

import (
	"net/http"

	"forestwhisperGo/magic/interfaces/http/middleware"
	"forestwhisperGo/magic/service"
	rbacSvc "forestwhisperGo/magic/service/rbac"

	"github.com/gin-gonic/gin"
)

func GetCurrentUserInfo(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusOK, middleware.Error("未登录或登录已过期喵喵", nil))
		return
	}

	profile := service.GetUserProfile(user.Uid)
	roles := rbacSvc.GetUserRoles(user.Uid)
	canAccessAdmin := rbacSvc.HasAnyAdminRole(user)

	data := gin.H{
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
	c.JSON(http.StatusOK, middleware.Success("ok", data))
}
