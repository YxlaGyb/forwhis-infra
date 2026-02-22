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
	var posts, pages, drafts, published, private, views int64
	db.GetDB().Model(&models.Content{}).Where("type = ?", models.TypePost).Count(&posts)
	db.GetDB().Model(&models.Content{}).Where("type = ?", models.TypePage).Count(&pages)
	db.GetDB().Model(&models.Content{}).Where("status = ?", models.StatusDraft).Count(&drafts)
	db.GetDB().Model(&models.Content{}).Where("status = ?", models.StatusPublish).Count(&published)
	db.GetDB().Model(&models.Content{}).Where("status = ?", models.StatusPrivate).Count(&private)
	db.GetDB().Model(&models.Content{}).Where("type = ?", models.TypePost).Select("COALESCE(SUM(views), 0)").Scan(&views)
	return map[string]int{"posts": int(posts), "pages": int(pages), "drafts": int(drafts), "published": int(published), "private": int(private), "views": int(views)}
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
	Views   int    `json:"views"`
}

type dailyCount struct {
	Date  string
	Count int
}

func GetTrendData(days int) []TrendData {
	if days <= 0 {
		days = 30
	}
	if days > 90 {
		days = 90
	}

	now := time.Now()
	loc := now.Location()
	startDate := now.AddDate(0, 0, -(days - 1))
	dayStart := int(time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc).Unix())

	dateMap := make(map[string]*TrendData)
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		dateMap[dateStr] = &TrendData{Date: dateStr}
	}

	var userCounts []dailyCount
	db.GetDB().Model(&models.User{}).
		Select("to_char(to_timestamp(created_at), 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("created_at >= ?", dayStart).
		Group("date").
		Scan(&userCounts)
	for _, uc := range userCounts {
		if td, ok := dateMap[uc.Date]; ok {
			td.User = uc.Count
		}
	}

	var contentCounts []dailyCount
	db.GetDB().Model(&models.Content{}).
		Select("to_char(to_timestamp(created), 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("created >= ?", dayStart).
		Group("date").
		Scan(&contentCounts)
	for _, cc := range contentCounts {
		if td, ok := dateMap[cc.Date]; ok {
			td.Content = cc.Count
		}
	}

	var commentCounts []dailyCount
	db.GetDB().Model(&models.Comment{}).
		Select("to_char(to_timestamp(created), 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("created >= ?", dayStart).
		Group("date").
		Scan(&commentCounts)
	for _, cc := range commentCounts {
		if td, ok := dateMap[cc.Date]; ok {
			td.Comment = cc.Count
		}
	}

	var viewsCounts []dailyCount
	db.GetDB().Model(&models.Content{}).
		Select("to_char(to_timestamp(created), 'YYYY-MM-DD') as date, COALESCE(SUM(views), 0) as count").
		Where("type = ? AND created >= ?", models.TypePost, dayStart).
		Group("date").
		Scan(&viewsCounts)
	for _, vc := range viewsCounts {
		if td, ok := dateMap[vc.Date]; ok {
			td.Views = vc.Count
		}
	}

	result := make([]TrendData, days)
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		result[i] = *dateMap[dateStr]
	}

	return result
}
