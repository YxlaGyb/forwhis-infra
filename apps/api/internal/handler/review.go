package handler

import (
	"strconv"

	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/request"
	"forestwhisperGo/apps/api/internal/service"
	rbacSvc "forestwhisperGo/apps/api/internal/service/rbac"
	"forestwhisperGo/apps/api/pkg/constants"

	"github.com/gofiber/fiber/v3"
)

func GetReviews(c fiber.Ctx) error {
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

	result := service.GetPendingReviews(page, limit)
	return middleware.JSONSuccess(c, "获取成功喵", result)
}

func ApproveReviewItem(c fiber.Ctx) error {
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return middleware.JSONError(c, "审核ID无效喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermCommentModerate) {
		return middleware.JSONError(c, "无权审核喵", nil)
	}

	success := service.ApproveReview(reviewId, user.Uid)
	if !success {
		return middleware.JSONError(c, "审核失败喵", nil)
	}
	return middleware.JSONSuccess(c, "审核通过喵", nil)
}

func RejectReviewItem(c fiber.Ctx) error {
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return middleware.JSONError(c, "审核ID无效喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermCommentModerate) {
		return middleware.JSONError(c, "无权审核喵", nil)
	}

	var req request.RejectReviewRequest
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请填写拒绝原因喵", nil)
	}

	success := service.RejectReview(reviewId, user.Uid, req.Reason)
	if !success {
		return middleware.JSONError(c, "审核失败喵", nil)
	}
	return middleware.JSONSuccess(c, "已拒绝喵", nil)
}
