package service

import (
	"regexp"
	"strconv"
	"time"

	"forestwhisperGo/apps/api/internal/models"
	"forestwhisperGo/apps/api/pkg/db"
)

func GetCategoriesList() []map[string]interface{} {
	var categories []models.Category
	db.GetDB().Order("\"order\" ASC, mid ASC").Find(&categories)
	counts := categoryContentCounts()
	items := []map[string]interface{}{}
	for i := range categories {
		items = append(items, categoryToDict(&categories[i], counts[categories[i].Mid]))
	}
	return items
}

func categoryContentCounts() map[int]int {
	var rows []struct {
		CategoryID int
		Count      int
	}
	db.GetDB().Model(&models.Content{}).
		Select("category_id, COUNT(*) AS count").
		Where("category_id > 0").
		Group("category_id").
		Scan(&rows)
	counts := make(map[int]int, len(rows))
	for _, row := range rows {
		counts[row.CategoryID] = row.Count
	}
	return counts
}

func GetCategoryByMid(mid int) *models.Category {
	var category models.Category
	if db.GetDB().First(&category, mid).Error != nil {
		return nil
	}
	return &category
}

func GetCategoryBySlug(slug string) *models.Category {
	var category models.Category
	if db.GetDB().Where("slug = ?", slug).First(&category).Error != nil {
		return nil
	}
	return &category
}

func CreateCategory(name, slug, description string, parent int) *models.Category {
	slug = uniqueSlug(slug, name, func(candidate string) bool {
		return GetCategoryBySlug(candidate) != nil
	})
	category := models.Category{
		Name: name, Slug: slug, Description: description, Parent: parent, Order: 0,
	}
	db.GetDB().Create(&category)
	return &category
}

func UpdateCategory(mid int, data map[string]interface{}) *models.Category {
	category := GetCategoryByMid(mid)
	if category == nil {
		return nil
	}
	updates := map[string]interface{}{}
	for _, k := range []string{"name", "slug", "description", "parent", "order"} {
		if v, ok := data[k]; ok {
			updates[k] = v
		}
	}
	db.GetDB().Model(category).Updates(updates)
	db.GetDB().First(category, mid)
	return category
}

func DeleteCategory(mid int) bool {
	category := GetCategoryByMid(mid)
	if category == nil {
		return false
	}
	db.GetDB().Model(&models.Content{}).Where("category_id = ?", mid).Update("category_id", 0)
	db.GetDB().Delete(category)
	return true
}

func categoryToDict(m *models.Category, count int) map[string]interface{} {
	return map[string]interface{}{
		"mid": m.Mid, "name": m.Name, "slug": m.Slug,
		"description": m.Description, "count": count, "parent": m.Parent,
		"type": "category",
	}
}

func GetCategoryCount(mid int) int {
	return categoryContentCounts()[mid]
}

func generateSlug(name string) string {
	slug := regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(name, "")
	slug = regexp.MustCompile(`[\s_-]+`).ReplaceAllString(slug, "-")
	return regexp.MustCompile(`^-|-$`).ReplaceAllString(slug, "")
}

func uniqueSlug(slug, name string, exists func(string) bool) string {
	if slug == "" {
		slug = generateSlug(name)
	}
	if !exists(slug) {
		return slug
	}
	return slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
}
