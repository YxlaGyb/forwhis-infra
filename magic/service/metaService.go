package service

import (
	"regexp"
	"time"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/utils/db"
)

func GetMetasList(metaType string) []map[string]interface{} {
	query := db.GetDB().Model(&models.Meta{})
	if metaType != "" {
		query = query.Where("type = ?", metaType)
	}
	var metas []models.Meta
	query.Order("\"order\" ASC, mid ASC").Find(&metas)
	items := []map[string]interface{}{}
	for _, m := range metas {
		items = append(items, metaToDict(&m))
	}
	return items
}

func GetMetaByMid(mid int) *models.Meta {
	var meta models.Meta
	if db.GetDB().First(&meta, mid).Error != nil {
		return nil
	}
	return &meta
}

func GetMetaBySlug(slug, metaType string) *models.Meta {
	var meta models.Meta
	if db.GetDB().Where("slug = ? AND type = ?", slug, metaType).First(&meta).Error != nil {
		return nil
	}
	return &meta
}

func CreateMeta(name, metaType, slug, description string, parent int) *models.Meta {
	if slug == "" {
		slug = generateMetaSlug(name)
	}
	if GetMetaBySlug(slug, metaType) != nil {
		slug = slug + "-" + string(rune(time.Now().Unix()))
	}
	meta := models.Meta{
		Name: name, Slug: slug, Type: metaType,
		Description: description, Parent: parent, Count: 0, Order: 0,
	}
	db.GetDB().Create(&meta)
	return &meta
}

func UpdateMeta(mid int, data map[string]interface{}) *models.Meta {
	meta := GetMetaByMid(mid)
	if meta == nil {
		return nil
	}
	updates := map[string]interface{}{}
	for _, k := range []string{"name", "slug", "description", "parent", "order"} {
		if v, ok := data[k]; ok {
			updates[k] = v
		}
	}
	db.GetDB().Model(meta).Updates(updates)
	return meta
}

func DeleteMeta(mid int) bool {
	meta := GetMetaByMid(mid)
	if meta == nil {
		return false
	}
	db.GetDB().Delete(meta)
	return true
}

func IncrementMetaCount(mid, delta int) {
	meta := GetMetaByMid(mid)
	if meta != nil {
		newCount := meta.Count + delta
		if newCount < 0 {
			newCount = 0
		}
		db.GetDB().Model(meta).Update("count", newCount)
	}
}

func metaToDict(m *models.Meta) map[string]interface{} {
	return map[string]interface{}{
		"mid": m.Mid, "name": m.Name, "slug": m.Slug,
		"type": m.Type, "description": m.Description,
		"count": m.Count, "parent": m.Parent,
	}
}

func generateMetaSlug(name string) string {
	slug := regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(name, "")
	slug = regexp.MustCompile(`[\s_-]+`).ReplaceAllString(slug, "-")
	return regexp.MustCompile(`^-|-$`).ReplaceAllString(slug, "")
}
