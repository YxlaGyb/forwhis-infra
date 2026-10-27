package handler

import (
	"strconv"

	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/service"

	"github.com/gofiber/fiber/v3"
)

func GetDashboard(c fiber.Ctx) error {
	data := service.GetDashboardData()
	return middleware.JSONSuccess(c, "获取成功喵", data)
}

func GetUserStats(c fiber.Ctx) error {
	data := service.GetUserStats()
	return middleware.JSONSuccess(c, "获取成功喵", data)
}

func GetContentStats(c fiber.Ctx) error {
	data := service.GetContentStats()
	return middleware.JSONSuccess(c, "获取成功喵", data)
}

func GetCommentStats(c fiber.Ctx) error {
	data := service.GetCommentStats()
	return middleware.JSONSuccess(c, "获取成功喵", data)
}

func GetTrendStats(c fiber.Ctx) error {
	daysStr := c.Query("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}
	data := service.GetTrendData(days)
	return middleware.JSONSuccess(c, "获取成功喵", data)
}
