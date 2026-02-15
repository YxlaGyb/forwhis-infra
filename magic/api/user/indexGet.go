package api

import (
	"net/http"

	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/service"

	"github.com/gin-gonic/gin"
)

func GetCurrentUserInfo(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusOK, middleware.Error("未登录或登录已过期喵喵", nil))
		return
	}

	profile := service.GetUserProfile(user.Uid)
	data := gin.H{
		"uid":       user.Uid,
		"name":      user.Name,
		"email":     user.Mail,
		"lastip":    user.LastIp,
		"createdAt": user.CreatedAt,
		"lastLogin": user.LastLogin,
		"isActive":  user.IsActive,
	}
	for k, v := range profile {
		data[k] = v
	}
	c.JSON(http.StatusOK, middleware.Success("获取成功", data))
}
