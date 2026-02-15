package controller

import (
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/service"

	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	data := service.GetDashboardData()
	middleware.JSONSuccess(c, "获取成功喵", data)
}

func GetUserStats(c *gin.Context) {
	data := service.GetUserStats()
	middleware.JSONSuccess(c, "获取成功喵", data)
}

func GetContentStats(c *gin.Context) {
	data := service.GetContentStats()
	middleware.JSONSuccess(c, "获取成功喵", data)
}

func GetCommentStats(c *gin.Context) {
	data := service.GetCommentStats()
	middleware.JSONSuccess(c, "获取成功喵", data)
}
