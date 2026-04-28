package service

import (
	"time"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/pkg/db"
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
		query = query.Joins("JOIN relationship ON content.cid = relationship.cid").
			Where("relationship.mid = ?", categoryMid)
	}
	var total int64
	query.Count(&total)
	var contents []models.Content
	offset := (page - 1) * limit
	query.Order("created desc").Offset(offset).Limit(limit).Find(&contents)

	categoryMap := batchGetContentCategories(contents)
	authorMap := batchGetContentAuthors(contents)

	items := []map[string]interface{}{}
	for _, c := range contents {
		items = append(items, contentToDictOptimized(&c, categoryMap, authorMap))
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

func batchGetContentCategories(contents []models.Content) map[int]*models.Meta {
	if len(contents) == 0 {
		return nil
	}
	cids := make([]int, len(contents))
	for i, c := range contents {
		cids[i] = c.Cid
	}

	var relationships []models.Relationship
	db.GetDB().Where("cid IN ?", cids).Find(&relationships)
	if len(relationships) == 0 {
		return nil
	}

	mids := make([]int, len(relationships))
	cidToMid := make(map[int]int)
	for i, r := range relationships {
		mids[i] = r.Mid
		cidToMid[r.Cid] = r.Mid
	}

	var metas []models.Meta
	db.GetDB().Where("mid IN ?", mids).Find(&metas)
	metaMap := make(map[int]*models.Meta)
	for i := range metas {
		metaMap[metas[i].Mid] = &metas[i]
	}

	result := make(map[int]*models.Meta)
	for cid, mid := range cidToMid {
		if meta, ok := metaMap[mid]; ok {
			result[cid] = meta
		}
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
	db.GetDB().Where("cid = ?", cid).Delete(&models.Comment{})
	db.GetDB().Where("cid = ?", cid).Delete(&models.Relationship{})
	db.GetDB().Where("target_type = ? AND target_id = ?", "content", cid).Delete(&models.Like{})
	db.GetDB().Delete(content)
	return true
}

func GetContentsByAuthor(authorID int) []map[string]interface{} {
	var contents []models.Content
	db.GetDB().Where("author_id = ?", authorID).Find(&contents)

	categoryMap := batchGetContentCategories(contents)
	authorMap := batchGetContentAuthors(contents)

	items := []map[string]interface{}{}
	for _, c := range contents {
		items = append(items, contentToDictOptimized(&c, categoryMap, authorMap))
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

func contentToDictOptimized(c *models.Content, categoryMap map[int]*models.Meta, authorMap map[int]map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"cid": c.Cid, "title": c.Title,
		"created": c.Created, "modified": c.Modified,
		"type": c.Type, "status": c.Status, "author_id": c.AuthorID,
		"views": c.Views, "comments_num": c.CommentsNum,
		"text": c.Text,
	}
	if categoryMap != nil {
		if category, ok := categoryMap[c.Cid]; ok {
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
	return result
}
