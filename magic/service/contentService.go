package service

import (
	"time"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/utils/db"
)

func GetContentsList(page, limit int, status string, authorID int) map[string]interface{} {
	query := db.GetDB().Model(&models.Content{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if authorID > 0 {
		query = query.Where("author_id = ?", authorID)
	}
	var total int64
	query.Count(&total)
	var contents []models.Content
	offset := (page - 1) * limit
	query.Order("created desc").Offset(offset).Limit(limit).Find(&contents)
	items := []map[string]interface{}{}
	for _, c := range contents {
		items = append(items, contentToDict(&c))
	}
	return map[string]interface{}{"total": total, "page": page, "limit": limit, "items": items}
}

func GetContentByCid(cid int) *models.Content {
	var content models.Content
	if db.GetDB().First(&content, cid).Error != nil {
		return nil
	}
	return &content
}

func IncrementViews(cid int) {
	db.GetDB().Model(&models.Content{}).Where("cid = ?", cid).UpdateColumn("views", db.GetDB().Raw("views + 1"))
}

func GetContentCategory(cid int) *models.Meta {
	var relationship models.Relationship
	if db.GetDB().Where("cid = ?", cid).First(&relationship).Error != nil {
		return nil
	}
	var meta models.Meta
	if db.GetDB().First(&meta, relationship.Mid).Error != nil {
		return nil
	}
	return &meta
}

func GetContentAuthor(authorID int) map[string]interface{} {
	var user models.User
	if db.GetDB().First(&user, authorID).Error != nil {
		return nil
	}
	return map[string]interface{}{
		"uid":  user.Uid,
		"name": user.Name,
	}
}

func CreateContent(title, text string, authorID int, contentType, status, password, template string, categoryMid int) *models.Content {
	now := int(time.Now().Unix())
	content := models.Content{
		Title: title, Text: text, AuthorID: authorID,
		Type: contentType, Status: status, Password: password, Template: template,
		Created: now, Modified: now,
	}
	if err := db.GetDB().Create(&content).Error; err != nil {
		return nil
	}

	if categoryMid > 0 {
		relationship := models.Relationship{Cid: content.Cid, Mid: categoryMid}
		db.GetDB().Create(&relationship)
	}

	return &content
}

func UpdateContent(cid int, data map[string]interface{}) *models.Content {
	content := GetContentByCid(cid)
	if content == nil {
		return nil
	}
	updates := map[string]interface{}{}
	for _, k := range []string{"title", "text", "status", "password", "template"} {
		if v, ok := data[k]; ok {
			updates[k] = v
		}
	}
	updates["modified"] = int(time.Now().Unix())
	db.GetDB().Model(content).Updates(updates)

	if categoryMid, ok := data["categoryMid"].(float64); ok {
		db.GetDB().Where("cid = ?", cid).Delete(&models.Relationship{})
		if int(categoryMid) > 0 {
			relationship := models.Relationship{Cid: cid, Mid: int(categoryMid)}
			db.GetDB().Create(&relationship)
		}
	}

	return content
}

func DeleteContent(cid int) bool {
	content := GetContentByCid(cid)
	if content == nil {
		return false
	}
	db.GetDB().Delete(content)
	db.GetDB().Where("cid = ?", cid).Delete(&models.Relationship{})
	return true
}

func GetContentsByAuthor(authorID int) []map[string]interface{} {
	var contents []models.Content
	db.GetDB().Where("author_id = ?", authorID).Find(&contents)
	items := []map[string]interface{}{}
	for _, c := range contents {
		items = append(items, contentToDict(&c))
	}
	return items
}

func contentToDict(c *models.Content) map[string]interface{} {
	result := map[string]interface{}{
		"cid": c.Cid, "title": c.Title,
		"created": c.Created, "modified": c.Modified,
		"type": c.Type, "status": c.Status, "author_id": c.AuthorID,
		"views": c.Views, "comments_num": c.CommentsNum,
		"text": c.Text,
	}
	category := GetContentCategory(c.Cid)
	if category != nil {
		result["category"] = map[string]interface{}{
			"mid":  category.Mid,
			"name": category.Name,
		}
	}
	author := GetContentAuthor(c.AuthorID)
	if author != nil {
		result["author"] = author
	}
	return result
}
