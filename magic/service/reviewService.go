package service

import (
	"encoding/json"
	"forestwhisperGo/magic/models"
	"forestwhisperGo/pkg/db"
	"time"
)

func CreateReview(targetType string, targetId int, reviewerId int, status int) *models.Review {
	now := int(time.Now().Unix())
	review := models.Review{
		TargetType: targetType,
		TargetId:   targetId,
		ReviewerId: reviewerId,
		Status:     status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	db.GetDB().Create(&review)
	return &review
}

func GetReviewByTarget(targetType string, targetId int) *models.Review {
	var review models.Review
	if db.GetDB().Where("target_type = ? AND target_id = ?", targetType, targetId).First(&review).Error != nil {
		return nil
	}
	return &review
}

func GetPendingReviews(page, limit int) map[string]interface{} {
	var reviews []models.Review
	var total int64

	db.GetDB().Model(&models.Review{}).Where("status = ?", models.ReviewStatusPending).Count(&total)

	offset := (page - 1) * limit
	db.GetDB().Where("status = ?", models.ReviewStatusPending).Order("created_at desc").Offset(offset).Limit(limit).Find(&reviews)

	items := []map[string]interface{}{}
	for _, r := range reviews {
		reviewData := map[string]interface{}{
			"review_id":   r.ReviewId,
			"target_type": r.TargetType,
			"target_id":   r.TargetId,
			"status":      r.Status,
			"reason":      r.Reason,
			"created":     r.CreatedAt,
		}

		if r.TargetType == models.TargetTypeUserProfile && r.Data != "" {
			var pendingData map[string]string
			json.Unmarshal([]byte(r.Data), &pendingData)
			reviewData["pending_data"] = pendingData
		}

		items = append(items, reviewData)
	}

	return map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
		"items": items,
	}
}

func CreateUserMetaReview(userId int, pendingFields map[string]string) {
	review := GetReviewByTarget(models.TargetTypeUserProfile, userId)
	if review != nil && review.Status == models.ReviewStatusPending {
		return
	}

	if review != nil {
		db.GetDB().Where("target_type = ? AND target_id = ?", models.TargetTypeUserProfile, userId).Delete(&models.Review{})
	}

	dataJSON, _ := json.Marshal(pendingFields)

	now := int(time.Now().Unix())
	review = &models.Review{
		TargetType: models.TargetTypeUserProfile,
		TargetId:   userId,
		ReviewerId: 0,
		Status:     models.ReviewStatusPending,
		Data:       string(dataJSON),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	db.GetDB().Create(review)
}

func ApproveReview(reviewId int, reviewerId int) bool {
	var review models.Review
	if db.GetDB().First(&review, reviewId).Error != nil {
		return false
	}

	if review.Status != models.ReviewStatusPending {
		return false
	}

	db.GetDB().Model(&review).Updates(map[string]interface{}{
		"status":      models.ReviewStatusApproved,
		"reviewer_id": reviewerId,
		"updated_at":  time.Now().Unix(),
	})

	applyReviewChanges(review.TargetType, review.TargetId)

	return true
}

func RejectReview(reviewId int, reviewerId int, reason string) bool {
	var review models.Review
	if db.GetDB().First(&review, reviewId).Error != nil {
		return false
	}

	if review.Status != models.ReviewStatusPending {
		return false
	}

	db.GetDB().Model(&review).Updates(map[string]interface{}{
		"status":      models.ReviewStatusRejected,
		"reviewer_id": reviewerId,
		"reason":      reason,
		"updated_at":  time.Now().Unix(),
	})

	return true
}

func applyReviewChanges(targetType string, targetId int) {
	switch targetType {
	case models.TargetTypeUserProfile:
		ApplyUserProfileReview(targetId)
	}
}

func ApplyUserProfileReview(userId int) {
	var review models.Review
	if db.GetDB().Where("target_type = ? AND target_id = ? AND status = ?",
		models.TargetTypeUserProfile, userId, models.ReviewStatusApproved).First(&review).Error != nil {
		return
	}

	if review.Data == "" {
		return
	}

	var pendingData map[string]string
	if err := json.Unmarshal([]byte(review.Data), &pendingData); err != nil {
		return
	}

	var user models.User
	if db.GetDB().First(&user, userId).Error != nil {
		return
	}

	if name, ok := pendingData["name"]; ok && name != "" {
		user.Name = name
	}

	db.GetDB().Save(&user)

	if bio, ok := pendingData["bio"]; ok {
		UpdateUsermeta(userId, models.MetaBio, bio)
	}
	if gender, ok := pendingData["gender"]; ok {
		UpdateUsermeta(userId, models.MetaGender, gender)
	}
	if birthday, ok := pendingData["birthday"]; ok {
		UpdateUsermeta(userId, models.MetaBirthday, birthday)
	}
}
