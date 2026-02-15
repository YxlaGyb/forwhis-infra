package controller

import (
	"net/http"

	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/service"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	email, password, clientIP := ctx.PostForm("email"), ctx.PostForm("password"), ctx.ClientIP()

	if !utils.IsValidEmail(email) || !utils.IsValidPassword(password) {
		ctx.JSON(http.StatusBadRequest, middleware.Error("邮箱或密码格式不正确喵", nil))
		return
	}

	cdKey := db.FormatKey("login:cd", email)
	if exists, _ := db.Exists(cdKey); exists > 0 {
		ctx.JSON(http.StatusTooManyRequests, middleware.Error("操作过于频繁，请稍后再试喵", nil))
		return
	}
	db.SetEX(cdKey, "1", 5)

	result, err := service.LoginUser(email, password, clientIP)
	if err != nil {
		db.Del(cdKey)
		ctx.JSON(http.StatusUnauthorized, middleware.Error(err.Error(), nil))
		return
	}

	token, _ := utils.GenerateToken(result.JID, result.UID, email)
	utils.SetCookieToken(ctx, token)
	ctx.JSON(http.StatusOK, middleware.Success("登录成功喵", gin.H{
		"uid": result.UID, "name": result.Name, "email": result.Email,
	}))
}
