package service

import (
	"time"

	"forestwhisperGo/apps/internal/models"
	"forestwhisperGo/apps/pkg/db"
)

func GetCommentsList(cid int, status string, authorID, page, limit int) map[string]interface{} {
	query := db.GetDB().Model(&models.Comment{})
	if cid > 0 {
		query = query.Where("cid = ?", cid)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if authorID > 0 {
		query = query.Where("author_id = ?", authorID)
	}
	var total int64
	query.Count(&total)
	var comments []models.Comment
	offset := (page - 1) * limit
	query.Order("created desc").Offset(offset).Limit(limit).Find(&comments)
	items := []map[string]interface{}{}
	for _, c := range comments {
		items = append(items, commentToDict(&c))
	}
	return map[string]interface{}{"total": total, "page": page, "limit": limit, "items": items}
}

func GetCommentByCoid(coid int) *models.Comment {
	var comment models.Comment
	if db.GetDB().First(&comment, coid).Error != nil {
		return nil
	}
	return &comment
}

func CreateComment(cid int, text string, authorID int, author, mail, url string, parent int) *models.Comment {
	comment := models.Comment{
		Cid: cid, Text: text, AuthorID: authorID, Author: author,
		Mail: mail, Url: url, Parent: parent,
		Created: int(time.Now().Unix()), Status: models.CommentStatusPending,
	}
	db.GetDB().Create(&comment)
	return &comment
}

func UpdateComment(coid int, data map[string]interface{}) *models.Comment {
	comment := GetCommentByCoid(coid)
	if comment == nil {
		return nil
	}
	updates := map[string]interface{}{}
	for _, k := range []string{"text", "status", "author", "mail", "url"} {
		if v, ok := data[k]; ok {
			updates[k] = v
		}
	}
	db.GetDB().Model(comment).Updates(updates)
	return comment
}

func ApproveComment(coid int) *models.Comment {
	return UpdateComment(coid, map[string]interface{}{"status": models.CommentStatusApproved})
}

func DeleteComment(coid int) bool {
	comment := GetCommentByCoid(coid)
	if comment == nil {
		return false
	}
	db.GetDB().Delete(comment)
	return true
}

func GetCommentsByContent(cid int, status string, viewerID int) []map[string]interface{} {
	query := db.GetDB().Where("cid = ?", cid)
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ? OR (status = ? AND author_id = ?)", models.CommentStatusApproved, models.CommentStatusPending, viewerID)
	}
	var comments []models.Comment
	query.Order("created asc").Find(&comments)
	items := []map[string]interface{}{}
	for _, c := range comments {
		items = append(items, commentToDict(&c))
	}
	return items
}

func GetAllCommentsByContent(cid int) []map[string]interface{} {
	query := db.GetDB().Where("cid = ?", cid)
	var comments []models.Comment
	query.Order("created asc").Find(&comments)
	items := []map[string]interface{}{}
	for _, c := range comments {
		items = append(items, commentToDict(&c))
	}
	return items
}

func GetCommentsByAuthor(authorID int) []map[string]interface{} {
	var comments []models.Comment
	db.GetDB().Where("author_id = ?", authorID).Order("created desc").Find(&comments)
	items := []map[string]interface{}{}
	for _, c := range comments {
		items = append(items, commentToDict(&c))
	}
	return items
}

func commentToDict(c *models.Comment) map[string]interface{} {
	return map[string]interface{}{
		"coid": c.Coid, "cid": c.Cid, "created": c.Created,
		"author": c.Author, "authorId": c.AuthorID, "text": c.Text,
		"status": c.Status, "parent": c.Parent,
		"mail": c.Mail, "url": c.Url,
	}
}
