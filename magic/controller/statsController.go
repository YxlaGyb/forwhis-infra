package controller

import (
	"strconv"

	"forestwhisperGo/magic/interfaces/http/middleware"
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

func GetTrendStats(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}
	data := service.GetTrendData(days)
	middleware.JSONSuccess(c, "获取成功喵", data)
}
