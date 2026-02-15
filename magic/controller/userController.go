package controller

import (
	"net/http"

	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/service"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	email, password, clientIP := c.PostForm("email"), c.PostForm("password"), c.ClientIP()

	if !utils.IsValidEmail(email) || !utils.IsValidPassword(password) {
		c.JSON(http.StatusBadRequest, middleware.Error("邮箱或密码格式不正确喵", nil))
		return
	}

	cdKey := db.FormatKey("login:cd", email)
	if exists, _ := db.Exists(cdKey); exists > 0 {
		c.JSON(http.StatusTooManyRequests, middleware.Error("操作过于频繁，请稍后再试喵", nil))
		return
	}
	db.SetEX(cdKey, "1", 5)

	result, errCode := service.LoginUser(email, password, clientIP)
	if errCode != 0 {
		db.Del(cdKey)
		c.JSON(http.StatusUnauthorized, middleware.Error("邮箱或密码错误喵", nil))
		return
	}

	token := result["token"].(string)
	utils.SetCookieToken(c, token)
	c.JSON(http.StatusOK, middleware.Success("登录成功喵", gin.H{
		"uid": result["uid"], "name": result["name"], "email": result["email"],
	}))
}
