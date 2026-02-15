package service

import (
	"regexp"
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

func GetContentBySlug(slug string) *models.Content {
	var content models.Content
	if db.GetDB().Where("slug = ?", slug).First(&content).Error != nil {
		return nil
	}
	return &content
}

func CreateContent(title, text string, authorID int, slug, contentType, status, password, template string) *models.Content {
	if slug == "" {
		slug = generateSlug(title)
	}
	if GetContentBySlug(slug) != nil {
		slug = slug + "-" + string(rune(time.Now().Unix()))
	}
	now := int(time.Now().Unix())
	content := models.Content{
		Title: title, Slug: slug, Text: text, AuthorID: authorID,
		Type: contentType, Status: status, Password: password, Template: template,
		Created: now, Modified: now,
	}
	db.GetDB().Create(&content)
	return &content
}

func UpdateContent(cid int, data map[string]interface{}) *models.Content {
	content := GetContentByCid(cid)
	if content == nil {
		return nil
	}
	updates := map[string]interface{}{}
	for _, k := range []string{"title", "text", "slug", "status", "password", "template"} {
		if v, ok := data[k]; ok {
			updates[k] = v
		}
	}
	updates["modified"] = int(time.Now().Unix())
	db.GetDB().Model(content).Updates(updates)
	return content
}

func DeleteContent(cid int) bool {
	content := GetContentByCid(cid)
	if content == nil {
		return false
	}
	db.GetDB().Delete(content)
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
	return map[string]interface{}{
		"cid": c.Cid, "title": c.Title, "slug": c.Slug,
		"created": c.Created, "modified": c.Modified,
		"type": c.Type, "status": c.Status, "author_id": c.AuthorID,
	}
}

func generateSlug(title string) string {
	slug := regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(title, "")
	slug = regexp.MustCompile(`[\s_-]+`).ReplaceAllString(slug, "-")
	return regexp.MustCompile(`^-|-$`).ReplaceAllString(slug, "")
}
