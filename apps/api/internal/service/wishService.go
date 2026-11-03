package service

import (
	"forwhis-infra/apps/api/internal/models"
	"forwhis-infra/apps/api/pkg/db"
	"time"
)

func GetWishes(page, limit int, viewerID int) map[string]interface{} {
	q := db.GetDB().Model(&models.Wish{}).Where("status = ?", models.WishStatusApproved)
	var total int64
	q.Count(&total)

	var wishes []models.Wish
	q.Order("created desc").Offset((page - 1) * limit).Limit(limit).Find(&wishes)

	items := []map[string]interface{}{}
	for i := range wishes {
		items = append(items, map[string]interface{}{
			"id": wishes[i].ID, "created": wishes[i].Created,
			"author": wishes[i].Author, "authorId": wishes[i].AuthorID,
			"content": wishes[i].Content,
			"likes":   GetLikeCount(models.LikeTargetTypeWish, wishes[i].ID),
			"liked":   viewerID > 0 && HasLiked(viewerID, models.LikeTargetTypeWish, wishes[i].ID),
		})
	}
	return map[string]interface{}{"total": total, "page": page, "limit": limit, "items": items}
}

func GetWishesByAuthor(authorID int) []map[string]interface{} {
	var wishes []models.Wish
	db.GetDB().Where("author_id = ?", authorID).Order("created desc").Find(&wishes)

	items := []map[string]interface{}{}
	for i := range wishes {
		items = append(items, map[string]interface{}{
			"id": wishes[i].ID, "created": wishes[i].Created,
			"author": wishes[i].Author, "authorId": wishes[i].AuthorID,
			"content": wishes[i].Content, "status": wishes[i].Status,
			"likes": GetLikeCount(models.LikeTargetTypeWish, wishes[i].ID),
			"liked": HasLiked(authorID, models.LikeTargetTypeWish, wishes[i].ID),
		})
	}
	return items
}

func DeleteWish(id, authorID int) bool {
	result := db.GetDB().Where("id = ? AND author_id = ?", id, authorID).Delete(&models.Wish{})
	if result.RowsAffected == 0 {
		return false
	}
	db.GetDB().Where("target_type = ? AND target_id = ?", models.LikeTargetTypeWish, id).Delete(&models.Like{})
	return true
}

func CreateWish(content string, authorID int, author, mail string) *models.Wish {
	wish := models.Wish{
		Content: content, AuthorID: authorID, Author: author,
		Mail: mail, Created: int(time.Now().Unix()), Status: models.WishStatusApproved,
	}
	db.GetDB().Create(&wish)
	return &wish
}
