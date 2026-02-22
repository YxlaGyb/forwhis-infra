package controller

import (
	"strconv"

	"forestwhisperGo/magic/constants"
	"forestwhisperGo/magic/middleware"
	rbacSvc "forestwhisperGo/magic/service/rbac"
	"forestwhisperGo/magic/service"

	"github.com/gin-gonic/gin"
)

func GetReviews(c *gin.Context) {
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

	result := service.GetPendingReviews(page, limit)
	middleware.JSONSuccess(c, "获取成功喵", result)
}

func ApproveReviewItem(c *gin.Context) {
	reviewId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		middleware.JSONError(c, "审核ID无效喵", nil)
		return
	}

	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermCommentModerate) {
		middleware.JSONError(c, "无权审核喵", nil)
		return
	}

	success := service.ApproveReview(reviewId, user.Uid)
	if !success {
		middleware.JSONError(c, "审核失败喵", nil)
		return
	}

	middleware.JSONSuccess(c, "审核通过喵", nil)
}

func RejectReviewItem(c *gin.Context) {
	reviewId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		middleware.JSONError(c, "审核ID无效喵", nil)
		return
	}

	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermCommentModerate) {
		middleware.JSONError(c, "无权审核喵", nil)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请填写拒绝原因喵", nil)
		return
	}

	success := service.RejectReview(reviewId, user.Uid, req.Reason)
	if !success {
		middleware.JSONError(c, "审核失败喵", nil)
		return
	}

	middleware.JSONSuccess(c, "已拒绝喵", nil)
}
