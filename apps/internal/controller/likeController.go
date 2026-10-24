package controller

import (
	"strconv"

	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/request"
	"forestwhisperGo/apps/internal/service"

	"github.com/gofiber/fiber/v3"
)

func ToggleLike(c fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return middleware.JSONError(c, "请先登录喵", nil)
	}

	var req request.ToggleLikeRequest

	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求参数无效喵", nil)
	}

	liked, err := service.ToggleLike(user.Uid, req.TargetType, req.TargetID)
	if err != nil {
		return middleware.JSONError(c, "操作失败喵", nil)
	}

	action := "取消点赞"
	if liked {
		action = "点赞"
	}

	service.CreateUserinfo(user.Uid, "toggleLike", action+"了"+req.TargetType+":"+strconv.Itoa(req.TargetID)+"喵")

	count := service.GetLikeCount(req.TargetType, req.TargetID)
	return middleware.JSONSuccess(c, "操作成功喵", map[string]interface{}{
		"liked": liked,
		"count": count,
	})

}

func GetLikeCount(c fiber.Ctx) error {
	targetType := c.Query("target_type")
	targetID, err := strconv.Atoi(c.Query("target_id"))
	if err != nil || targetID <= 0 {
		return middleware.JSONError(c, "目标ID无效喵", nil)
	}

	count := service.GetLikeCount(targetType, targetID)
	return middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
		"count": count,
	})

}

func GetLikeStatus(c fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	targetType := c.Query("target_type")
	targetID, err := strconv.Atoi(c.Query("target_id"))
	if err != nil || targetID <= 0 {
		return middleware.JSONError(c, "目标ID无效喵", nil)
	}

	var hasLiked bool
	if user != nil {
		hasLiked = service.HasLiked(user.Uid, targetType, targetID)
	}

	count := service.GetLikeCount(targetType, targetID)
	return middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
		"has_liked": hasLiked,
		"count":     count,
	})

}

func GetUserLikes(c fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return middleware.JSONError(c, "请先登录喵", nil)
	}

	targetType := c.Query("target_type")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	result := service.GetLikesByUser(user.Uid, targetType, page, limit)
	return middleware.JSONSuccess(c, "获取成功喵", result)
}

func GetMyLikesCount(c fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return middleware.JSONError(c, "请先登录喵", nil)
	}

	count := service.GetUserLikedCount(user.Uid)
	return middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
		"count": count,
	})

}
