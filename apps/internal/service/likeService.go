package service

import (
	"time"

	"forestwhisperGo/apps/internal/models"
	"forestwhisperGo/apps/pkg/db"
)

func ToggleLike(userID int, targetType string, targetID int) (bool, error) {
	var existingLike models.Like
	result := db.GetDB().Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).First(&existingLike)

	if result.Error == nil {
		db.GetDB().Delete(&existingLike)
		return false, nil
	}

	like := models.Like{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   targetID,
		Created:    int(time.Now().Unix()),
	}

	if err := db.GetDB().Create(&like).Error; err != nil {
		return false, err
	}

	return true, nil
}

func GetLikeCount(targetType string, targetID int) int64 {
	var count int64
	db.GetDB().Model(&models.Like{}).Where("target_type = ? AND target_id = ?", targetType, targetID).Count(&count)
	return count
}

func HasLiked(userID int, targetType string, targetID int) bool {
	var count int64
	db.GetDB().Model(&models.Like{}).Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Count(&count)
	return count > 0
}

func GetLikesByUser(userID int, targetType string, page, limit int) map[string]interface{} {
	query := db.GetDB().Where("user_id = ?", userID)
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	var total int64
	query.Model(&models.Like{}).Count(&total)

	var likes []models.Like
	offset := (page - 1) * limit
	query.Order("created desc").Offset(offset).Limit(limit).Find(&likes)

	items := []map[string]interface{}{}
	for _, like := range likes {
		items = append(items, map[string]interface{}{
			"id":          like.ID,
			"target_type": like.TargetType,
			"target_id":   like.TargetID,
			"created":     like.Created,
		})
	}

	return map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
		"items": items,
	}
}

func GetLikeStatus(userID int, targetType string, targetID int) map[string]interface{} {
	hasLiked := HasLiked(userID, targetType, targetID)
	count := GetLikeCount(targetType, targetID)

	return map[string]interface{}{
		"has_liked": hasLiked,
		"count":     count,
	}
}

func GetUserLikedCount(userID int) int64 {
	var count int64
	db.GetDB().Model(&models.Like{}).Where("user_id = ?", userID).Count(&count)
	return count
}
