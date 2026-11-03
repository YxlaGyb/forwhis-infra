package feed

import (
	"fmt"

	"forwhis-infra/apps/api/internal/models"
	"forwhis-infra/apps/api/pkg/db"

	"gorm.io/gorm"
)

// ponytail: fixed recall cap (200 per source); deep paging stops around 400.
const recallLimit = 200

func Recall(snapshot int64) ([]Candidate, error) {
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("feed: database unavailable")
	}

	var contents []models.Content
	if err := database.
		Where("status = ? AND type = ? AND password = ? AND created <= ?",
			models.StatusPublish, models.TypePost, "", snapshot).
		Order("created desc").
		Limit(recallLimit).
		Find(&contents).Error; err != nil {
		return nil, err
	}

	var wishes []models.Wish
	if err := database.
		Where("status = ? AND created <= ?", models.WishStatusApproved, snapshot).
		Order("created desc").
		Limit(recallLimit).
		Find(&wishes).Error; err != nil {
		return nil, err
	}

	cids := make([]int, len(contents))
	authorIDs := make([]int, 0, len(contents))
	categoryIDs := make([]int, 0, len(contents))
	seenAuthors := make(map[int]struct{}, len(contents))
	seenCategories := make(map[int]struct{}, len(contents))
	for i, c := range contents {
		cids[i] = c.Cid
		if _, ok := seenAuthors[c.AuthorID]; !ok && c.AuthorID > 0 {
			seenAuthors[c.AuthorID] = struct{}{}
			authorIDs = append(authorIDs, c.AuthorID)
		}
		if _, ok := seenCategories[c.CategoryID]; !ok && c.CategoryID > 0 {
			seenCategories[c.CategoryID] = struct{}{}
			categoryIDs = append(categoryIDs, c.CategoryID)
		}
	}

	wishIDs := make([]int, len(wishes))
	for i, w := range wishes {
		wishIDs[i] = w.ID
	}

	authorNames := loadAuthorNames(database, authorIDs)
	categoryNames := loadCategoryNames(database, categoryIDs)
	tags := loadTags(database, cids)
	commentCounts := loadCommentCounts(database, cids)
	contentLikes := loadLikeCounts(database, models.LikeTargetTypeContent, cids)
	wishLikes := loadLikeCounts(database, models.LikeTargetTypeWish, wishIDs)

	candidates := make([]Candidate, 0, len(contents)+len(wishes))
	for _, c := range contents {
		candidates = append(candidates, Candidate{
			Kind:         KindContent,
			ID:           c.Cid,
			AuthorID:     c.AuthorID,
			AuthorName:   authorNames[c.AuthorID],
			Title:        c.Title,
			Text:         c.Text,
			CategoryID:   c.CategoryID,
			CategoryName: categoryNames[c.CategoryID],
			Tags:         tags[c.Cid],
			Created:      int64(c.Created),
			Views:        c.Views,
			Comments:     commentCounts[c.Cid],
			Likes:        contentLikes[c.Cid],
			Status:       c.Status,
			Type:         c.Type,
			HasPassword:  c.Password != "",
		})
	}
	for _, w := range wishes {
		candidates = append(candidates, Candidate{
			Kind:       KindWish,
			ID:         w.ID,
			AuthorID:   w.AuthorID,
			AuthorName: w.Author,
			Text:       w.Content,
			Created:    int64(w.Created),
			Likes:      wishLikes[w.ID],
			Status:     w.Status,
		})
	}
	return candidates, nil
}

func loadAuthorNames(database *gorm.DB, ids []int) map[int]string {
	result := make(map[int]string, len(ids))
	if len(ids) == 0 {
		return result
	}

	var users []models.User
	database.Select("uid, name").Where("uid IN ?", ids).Find(&users)
	for _, u := range users {
		result[u.Uid] = u.Name
	}
	return result
}

func loadCategoryNames(database *gorm.DB, ids []int) map[int]string {
	result := make(map[int]string, len(ids))
	if len(ids) == 0 {
		return result
	}

	var categories []models.Category
	database.Select("mid, name").Where("mid IN ?", ids).Find(&categories)
	for _, c := range categories {
		result[c.Mid] = c.Name
	}
	return result
}

func loadTags(database *gorm.DB, cids []int) map[int][]string {
	result := make(map[int][]string)
	if len(cids) == 0 {
		return result
	}

	var rows []struct {
		Cid  int
		Name string
	}
	database.Table("content_tags").
		Select("content_tags.cid, tags.name").
		Joins("JOIN tags ON tags.id = content_tags.tag_id").
		Where("content_tags.cid IN ?", cids).
		Order("tags.name ASC").
		Scan(&rows)

	for _, row := range rows {
		result[row.Cid] = append(result[row.Cid], row.Name)
	}
	return result
}

func loadCommentCounts(database *gorm.DB, cids []int) map[int]int {
	result := make(map[int]int)
	if len(cids) == 0 {
		return result
	}

	var rows []struct {
		Cid   int
		Count int
	}
	database.Model(&models.Comment{}).
		Select("cid, COUNT(*) AS count").
		Where("cid IN ?", cids).
		Group("cid").
		Scan(&rows)

	for _, row := range rows {
		result[row.Cid] = row.Count
	}
	return result
}

func loadLikeCounts(database *gorm.DB, targetType string, ids []int) map[int]int {
	result := make(map[int]int)
	if len(ids) == 0 {
		return result
	}

	var rows []struct {
		TargetID int
		Count    int
	}
	database.Model(&models.Like{}).
		Select("target_id, COUNT(*) AS count").
		Where("target_type = ? AND target_id IN ?", targetType, ids).
		Group("target_id").
		Scan(&rows)

	for _, row := range rows {
		result[row.TargetID] = row.Count
	}
	return result
}
