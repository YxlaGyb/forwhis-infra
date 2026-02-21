package controller

import (
	"strconv"

	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/service"

	"github.com/gin-gonic/gin"
)

func ToggleLike(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		middleware.JSONError(c, "请先登录喵", nil)
		return
	}

	var req struct {
		TargetType string `json:"target_type" binding:"required,oneof=content comment"`
		TargetID   int    `json:"target_id" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求参数无效喵", nil)
		return
	}

	liked, err := service.ToggleLike(user.Uid, req.TargetType, req.TargetID)
	if err != nil {
		middleware.JSONError(c, "操作失败喵", nil)
		return
	}

	action := "取消点赞"
	if liked {
		action = "点赞"
	}

	service.CreateUserinfo(user.Uid, "toggleLike", action+"了"+req.TargetType+":"+strconv.Itoa(req.TargetID)+"喵")

	count := service.GetLikeCount(req.TargetType, req.TargetID)
	middleware.JSONSuccess(c, "操作成功喵", map[string]interface{}{
		"liked": liked,
		"count": count,
	})
}

func GetLikeCount(c *gin.Context) {
	targetType := c.Query("target_type")
	targetID, err := strconv.Atoi(c.Query("target_id"))
	if err != nil || targetID <= 0 {
		middleware.JSONError(c, "目标ID无效喵", nil)
		return
	}

	count := service.GetLikeCount(targetType, targetID)
	middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
		"count": count,
	})
}

func GetLikeStatus(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	targetType := c.Query("target_type")
	targetID, err := strconv.Atoi(c.Query("target_id"))
	if err != nil || targetID <= 0 {
		middleware.JSONError(c, "目标ID无效喵", nil)
		return
	}

	var hasLiked bool
	if user != nil {
		hasLiked = service.HasLiked(user.Uid, targetType, targetID)
	}

	count := service.GetLikeCount(targetType, targetID)
	middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
		"has_liked": hasLiked,
		"count":     count,
	})
}

func GetUserLikes(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		middleware.JSONError(c, "请先登录喵", nil)
		return
	}

	targetType := c.Query("target_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

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
	middleware.JSONSuccess(c, "获取成功喵", result)
}
