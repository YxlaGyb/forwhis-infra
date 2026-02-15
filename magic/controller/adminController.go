package controller

import (
	"strconv"
	"time"

	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/service"
	"forestwhisperGo/magic/utils/db"

	"github.com/gin-gonic/gin"
)

func GetPendingReviews(c *gin.Context) {
	var reviews []models.Review
	db.GetDB().Where("status = ?", models.ReviewStatusPending).Find(&reviews)

	data := []map[string]interface{}{}
	for _, r := range reviews {
		data = append(data, map[string]interface{}{
			"review_id":   r.ReviewId,
			"target_type": r.TargetType,
			"target_id":   r.TargetId,
			"created_at":  r.CreatedAt,
		})
	}
	middleware.JSONSuccess(c, "获取成功喵", data)
}

func ApproveReview(c *gin.Context) {
	reviewId, err := strconv.Atoi(c.Param("review_id"))
	if err != nil {
		middleware.JSONError(c, "审核ID无效喵", nil)
		return
	}

	var req struct {
		NewValue string `json:"newValue"`
	}
	c.ShouldBindJSON(&req)

	var review models.Review
	if db.GetDB().First(&review, reviewId).Error != nil {
		middleware.JSONError(c, "审核记录不存在喵", nil)
		return
	}

	if review.TargetType == models.TargetTypeUserProfile {
		var meta models.Usermeta
		if db.GetDB().First(&meta, review.TargetId).Error == nil && req.NewValue != "" {
			meta.MetaValue = req.NewValue
			meta.UpdatedAt = int(time.Now().Unix())
			db.GetDB().Save(&meta)
		}
	}

	review.Status = models.ReviewStatusApproved
	review.UpdatedAt = int(time.Now().Unix())
	user := middleware.GetCurrentUser(c)
	if user != nil {
		review.ReviewerId = user.Uid
	}
	db.GetDB().Save(&review)

	if user != nil {
		service.CreateUserinfo(user.Uid, "approveReview", "通过了审核: "+strconv.Itoa(reviewId)+"喵")
	}

	middleware.JSONSuccess(c, "审核已通过喵", nil)
}

func RejectReview(c *gin.Context) {
	reviewId, err := strconv.Atoi(c.Param("review_id"))
	if err != nil {
		middleware.JSONError(c, "审核ID无效喵", nil)
		return
	}

	var review models.Review
	if db.GetDB().First(&review, reviewId).Error != nil {
		middleware.JSONError(c, "审核记录不存在喵", nil)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	review.Status = models.ReviewStatusRejected
	review.UpdatedAt = int(time.Now().Unix())
	review.Reason = req.Reason
	user := middleware.GetCurrentUser(c)
	if user != nil {
		review.ReviewerId = user.Uid
	}
	db.GetDB().Save(&review)

	if user != nil {
		service.CreateUserinfo(user.Uid, "rejectReview", "拒绝了审核: "+strconv.Itoa(reviewId)+"喵")
	}

	middleware.JSONSuccess(c, "审核已拒绝喵", nil)
}
