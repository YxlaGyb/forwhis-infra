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
	db.GetDB().Model(&models.User{}).Count(&total)
	db.GetDB().Model(&models.User{}).Where("created_at >= ?", todayStart).Count(&todayNew)
	db.GetDB().Model(&models.User{}).Where("created_at >= ?", weekStart).Count(&weekNew)
	return map[string]int{"total": int(total), "todayNew": int(todayNew), "weekNew": int(weekNew)}
}

func GetContentStats() map[string]int {
	var posts, pages, drafts, published, private int64
	db.GetDB().Model(&models.Content{}).Where("type = ?", models.TypePost).Count(&posts)
	db.GetDB().Model(&models.Content{}).Where("type = ?", models.TypePage).Count(&pages)
	db.GetDB().Model(&models.Content{}).Where("status = ?", models.StatusDraft).Count(&drafts)
	db.GetDB().Model(&models.Content{}).Where("status = ?", models.StatusPublish).Count(&published)
	db.GetDB().Model(&models.Content{}).Where("status = ?", models.StatusPrivate).Count(&private)
	return map[string]int{"posts": int(posts), "pages": int(pages), "drafts": int(drafts), "published": int(published), "private": int(private)}
}

func GetCommentStats() map[string]int {
	var total, pending, approved, spam int64
	db.GetDB().Model(&models.Comment{}).Count(&total)
	db.GetDB().Model(&models.Comment{}).Where("status = ?", models.CommentStatusPending).Count(&pending)
	db.GetDB().Model(&models.Comment{}).Where("status = ?", models.CommentStatusApproved).Count(&approved)
	db.GetDB().Model(&models.Comment{}).Where("status = ?", models.CommentStatusSpam).Count(&spam)
	return map[string]int{"total": int(total), "pending": int(pending), "approved": int(approved), "spam": int(spam)}
}

func GetDashboardData() map[string]map[string]int {
	return map[string]map[string]int{
		"users":    GetUserStats(),
		"contents": GetContentStats(),
		"comments": GetCommentStats(),
	}
}
