package service

import (
	"time"

	"forestwhisperGo/apps/internal/models"
	"forestwhisperGo/apps/pkg/db"

	"gorm.io/gorm"
)

func GetTagsList() []map[string]interface{} {
	var tags []models.Tag
	db.GetDB().Order("name ASC").Find(&tags)
	items := []map[string]interface{}{}
	for i := range tags {
		items = append(items, tagToDict(&tags[i]))
	}
	return items
}

func GetTagById(id int) *models.Tag {
	var tag models.Tag
	if db.GetDB().First(&tag, id).Error != nil {
		return nil
	}
	return &tag
}

func GetTagBySlug(slug string) *models.Tag {
	var tag models.Tag
	if db.GetDB().Where("slug = ?", slug).First(&tag).Error != nil {
		return nil
	}
	return &tag
}

func CreateTag(name, slug string) *models.Tag {
	slug = uniqueSlug(slug, name, func(candidate string) bool {
		return GetTagBySlug(candidate) != nil
	})
	tag := models.Tag{
		Name: name, Slug: slug, CreatedAt: int(time.Now().Unix()),
	}
	db.GetDB().Create(&tag)
	return &tag
}

func UpdateTag(id int, data map[string]interface{}) *models.Tag {
	tag := GetTagById(id)
	if tag == nil {
		return nil
	}
	updates := map[string]interface{}{}
	for _, k := range []string{"name", "slug"} {
		if v, ok := data[k]; ok {
			updates[k] = v
		}
	}
	db.GetDB().Model(tag).Updates(updates)
	db.GetDB().First(tag, id)
	return tag
}

func DeleteTag(id int) bool {
	tag := GetTagById(id)
	if tag == nil {
		return false
	}
	db.GetDB().Where("tag_id = ?", id).Delete(&models.ContentTag{})
	db.GetDB().Delete(tag)
	return true
}

func SetContentTags(cid int, tagIds []int) error {
	unique := make([]int, 0, len(tagIds))
	seen := make(map[int]struct{}, len(tagIds))
	for _, id := range tagIds {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}

	valid := make([]models.Tag, 0, len(unique))
	if len(unique) > 0 {
		db.GetDB().Where("id IN ?", unique).Find(&valid)
	}

	return db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("cid = ?", cid).Delete(&models.ContentTag{}).Error; err != nil {
			return err
		}
		for _, tag := range valid {
			if err := tx.Create(&models.ContentTag{Cid: cid, TagId: tag.Id}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func batchGetContentTags(cids []int) map[int][]models.Tag {
	result := make(map[int][]models.Tag)
	if len(cids) == 0 {
		return result
	}

	var rows []struct {
		Cid  int
		Id   int
		Name string
		Slug string
	}
	db.GetDB().Table("content_tags").
		Select("content_tags.cid, tags.id, tags.name, tags.slug").
		Joins("JOIN tags ON tags.id = content_tags.tag_id").
		Where("content_tags.cid IN ?", cids).
		Order("tags.name ASC").
		Scan(&rows)

	for _, row := range rows {
		result[row.Cid] = append(result[row.Cid], models.Tag{
			Id: row.Id, Name: row.Name, Slug: row.Slug,
		})
	}
	return result
}

func GetContentTags(cid int) []map[string]interface{} {
	tags := batchGetContentTags([]int{cid})[cid]
	items := []map[string]interface{}{}
	for i := range tags {
		items = append(items, tagToDict(&tags[i]))
	}
	return items
}

func tagToDict(t *models.Tag) map[string]interface{} {
	return map[string]interface{}{
		"id": t.Id, "mid": t.Id, "name": t.Name, "slug": t.Slug, "type": "tag",
	}
}
