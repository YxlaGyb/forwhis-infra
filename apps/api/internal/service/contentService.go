package service

import (
	"time"

	"forestwhisperGo/apps/api/internal/models"
	"forestwhisperGo/apps/api/pkg/db"

	"gorm.io/gorm"
)

func GetContentsList(page, limit int, status string, authorID, categoryMid int) map[string]interface{} {
	query := db.GetDB().Model(&models.Content{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if authorID > 0 {
		query = query.Where("author_id = ?", authorID)
	}
	if categoryMid > 0 {
		query = query.Where("category_id = ?", categoryMid)
	}
	var total int64
	query.Count(&total)
	var contents []models.Content
	offset := (page - 1) * limit
	query.Order("created desc").Offset(offset).Limit(limit).Find(&contents)

	categoryMap := batchGetContentCategories(contents)
	authorMap := batchGetContentAuthors(contents)
	tagMap := batchGetContentTagsFor(contents)
	commentCounts := batchGetCommentCounts(contents)

	items := []map[string]interface{}{}
	for _, c := range contents {
		items = append(items, contentToDictOptimized(&c, categoryMap, authorMap, tagMap, commentCounts))
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
	db.GetDB().Model(&models.Content{}).Where("cid = ?", cid).UpdateColumn("views", gorm.Expr("views + 1"))
}

func GetContentCategory(cid int) *models.Category {
	content := GetContentByCid(cid)
	if content == nil || content.CategoryID == 0 {
		return nil
	}
	return GetCategoryByMid(content.CategoryID)
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

func GetCommentCount(cid int) int64 {
	var count int64
	db.GetDB().Model(&models.Comment{}).Where("cid = ?", cid).Count(&count)
	return count
}

func batchGetContentCategories(contents []models.Content) map[int]*models.Category {
	if len(contents) == 0 {
		return nil
	}
	ids := make([]int, 0, len(contents))
	seen := make(map[int]struct{}, len(contents))
	for _, content := range contents {
		if content.CategoryID == 0 {
			continue
		}
		if _, ok := seen[content.CategoryID]; ok {
			continue
		}
		seen[content.CategoryID] = struct{}{}
		ids = append(ids, content.CategoryID)
	}
	if len(ids) == 0 {
		return nil
	}

	var categories []models.Category
	db.GetDB().Where("mid IN ?", ids).Find(&categories)
	result := make(map[int]*models.Category, len(categories))
	for i := range categories {
		result[categories[i].Mid] = &categories[i]
	}
	return result
}

func batchGetContentAuthors(contents []models.Content) map[int]map[string]interface{} {
	if len(contents) == 0 {
		return nil
	}
	authorIDs := make(map[int]bool)
	for _, c := range contents {
		authorIDs[c.AuthorID] = true
	}
	ids := make([]int, 0, len(authorIDs))
	for id := range authorIDs {
		ids = append(ids, id)
	}

	var users []models.User
	db.GetDB().Where("uid IN ?", ids).Find(&users)

	result := make(map[int]map[string]interface{})
	for _, u := range users {
		result[u.Uid] = map[string]interface{}{
			"uid":  u.Uid,
			"name": u.Name,
		}
	}
	return result
}

func batchGetContentTagsFor(contents []models.Content) map[int][]models.Tag {
	cids := make([]int, len(contents))
	for i, content := range contents {
		cids[i] = content.Cid
	}
	return batchGetContentTags(cids)
}

func batchGetCommentCounts(contents []models.Content) map[int]int {
	if len(contents) == 0 {
		return nil
	}
	cids := make([]int, len(contents))
	for i, content := range contents {
		cids[i] = content.Cid
	}
	var rows []struct {
		Cid   int
		Count int
	}
	db.GetDB().Model(&models.Comment{}).
		Select("cid, COUNT(*) AS count").
		Where("cid IN ?", cids).
		Group("cid").
		Scan(&rows)
	result := make(map[int]int, len(rows))
	for _, row := range rows {
		result[row.Cid] = row.Count
	}
	return result
}

func CreateContent(title, text string, authorID int, contentType, status, password string, categoryMid int, tagIds []int) *models.Content {
	now := int(time.Now().Unix())
	content := models.Content{
		Title: title, Text: text, AuthorID: authorID,
		Type: contentType, Status: status, Password: password,
		CategoryID: categoryMid,
		Created:    now, Modified: now,
	}
	if err := db.GetDB().Create(&content).Error; err != nil {
		return nil
	}

	if len(tagIds) > 0 {
		if err := SetContentTags(content.Cid, tagIds); err != nil {
			db.GetDB().Delete(&content)
			return nil
		}
	}

	return &content
}

func UpdateContent(cid int, data map[string]interface{}) *models.Content {
	content := GetContentByCid(cid)
	if content == nil {
		return nil
	}
	updates := map[string]interface{}{}
	for _, k := range []string{"title", "text", "status", "password"} {
		if v, ok := data[k]; ok {
			updates[k] = v
		}
	}
	if categoryMid, ok := data["categoryMid"].(float64); ok {
		updates["category_id"] = int(categoryMid)
	}
	updates["modified"] = int(time.Now().Unix())
	db.GetDB().Model(content).Updates(updates)

	if raw, ok := data["tagIds"]; ok {
		SetContentTags(cid, parseIntSlice(raw))
	}

	db.GetDB().First(content, cid)
	return content
}

func parseIntSlice(value interface{}) []int {
	raw, ok := value.([]interface{})
	if !ok {
		return nil
	}
	ids := make([]int, 0, len(raw))
	for _, item := range raw {
		if number, ok := item.(float64); ok {
			ids = append(ids, int(number))
		}
	}
	return ids
}

func DeleteContent(cid int) bool {
	content := GetContentByCid(cid)
	if content == nil {
		return false
	}
	db.GetDB().Where("cid = ?", cid).Delete(&models.Comment{})
	db.GetDB().Where("cid = ?", cid).Delete(&models.ContentTag{})
	db.GetDB().Where("target_type = ? AND target_id = ?", "content", cid).Delete(&models.Like{})
	db.GetDB().Delete(content)
	return true
}

func GetContentsByAuthor(authorID int) []map[string]interface{} {
	var contents []models.Content
	db.GetDB().Where("author_id = ?", authorID).Find(&contents)

	categoryMap := batchGetContentCategories(contents)
	authorMap := batchGetContentAuthors(contents)
	tagMap := batchGetContentTagsFor(contents)
	commentCounts := batchGetCommentCounts(contents)

	items := []map[string]interface{}{}
	for _, c := range contents {
		items = append(items, contentToDictOptimized(&c, categoryMap, authorMap, tagMap, commentCounts))
	}
	return items
}

func contentToDictOptimized(c *models.Content, categoryMap map[int]*models.Category, authorMap map[int]map[string]interface{}, tagMap map[int][]models.Tag, commentCounts map[int]int) map[string]interface{} {
	result := map[string]interface{}{
		"cid": c.Cid, "title": c.Title,
		"created": c.Created, "modified": c.Modified,
		"type": c.Type, "status": c.Status, "author_id": c.AuthorID,
		"views": c.Views, "comments_num": commentCounts[c.Cid],
		"text": c.Text,
	}
	if categoryMap != nil {
		if category, ok := categoryMap[c.CategoryID]; ok {
			result["category"] = map[string]interface{}{
				"mid":  category.Mid,
				"name": category.Name,
			}
		}
	}
	if authorMap != nil {
		if author, ok := authorMap[c.AuthorID]; ok {
			result["author"] = author
		}
	}
	tags := []map[string]interface{}{}
	for i := range tagMap[c.Cid] {
		tags = append(tags, tagToDict(&tagMap[c.Cid][i]))
	}
	result["tags"] = tags
	return result
}
