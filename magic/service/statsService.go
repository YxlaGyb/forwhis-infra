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

type TrendData struct {
	Date    string `json:"date"`
	User    int    `json:"user"`
	Content int    `json:"content"`
	Comment int    `json:"comment"`
}

func GetTrendData(days int) []TrendData {
	if days <= 0 {
		days = 30
	}
	if days > 90 {
		days = 90
	}

	result := make([]TrendData, days)
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		dayStart := int(date.Unix())
		dayEnd := dayStart + 86400

		var userCount, contentCount, commentCount int64

		db.GetDB().Model(&models.User{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&userCount)
		db.GetDB().Model(&models.Content{}).Where("created >= ? AND created < ?", dayStart, dayEnd).Count(&contentCount)
		db.GetDB().Model(&models.Comment{}).Where("created >= ? AND created < ?", dayStart, dayEnd).Count(&commentCount)

		result[days-1-i] = TrendData{
			Date:    dateStr,
			User:    int(userCount),
			Content: int(contentCount),
			Comment: int(commentCount),
		}
	}

	return result
}
