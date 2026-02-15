package service

import (
	"time"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/utils/db"
)

func GetUserStats() map[string]int {
	now := int(time.Now().Unix())
	todayStart := now - (now % 86400)
	weekStart := now - (7 * 86400)
	var total, todayNew, weekNew int64
	db.DB.Model(&models.User{}).Count(&total)
	db.DB.Model(&models.User{}).Where("created_at >= ?", todayStart).Count(&todayNew)
	db.DB.Model(&models.User{}).Where("created_at >= ?", weekStart).Count(&weekNew)
	return map[string]int{"total": int(total), "todayNew": int(todayNew), "weekNew": int(weekNew)}
}

func GetContentStats() map[string]int {
	var posts, pages, drafts, published, private int64
	db.DB.Model(&models.Content{}).Where("type = ?", models.TypePost).Count(&posts)
	db.DB.Model(&models.Content{}).Where("type = ?", models.TypePage).Count(&pages)
	db.DB.Model(&models.Content{}).Where("status = ?", models.StatusDraft).Count(&drafts)
	db.DB.Model(&models.Content{}).Where("status = ?", models.StatusPublish).Count(&published)
	db.DB.Model(&models.Content{}).Where("status = ?", models.StatusPrivate).Count(&private)
	return map[string]int{"posts": int(posts), "pages": int(pages), "drafts": int(drafts), "published": int(published), "private": int(private)}
}

func GetCommentStats() map[string]int {
	var total, pending, approved, spam int64
	db.DB.Model(&models.Comment{}).Count(&total)
	db.DB.Model(&models.Comment{}).Where("status = ?", models.CommentStatusPending).Count(&pending)
	db.DB.Model(&models.Comment{}).Where("status = ?", models.CommentStatusApproved).Count(&approved)
	db.DB.Model(&models.Comment{}).Where("status = ?", models.CommentStatusSpam).Count(&spam)
	return map[string]int{"total": int(total), "pending": int(pending), "approved": int(approved), "spam": int(spam)}
}

func GetDashboardData() map[string]map[string]int {
	return map[string]map[string]int{
		"users":    GetUserStats(),
		"contents": GetContentStats(),
		"comments": GetCommentStats(),
	}
}
