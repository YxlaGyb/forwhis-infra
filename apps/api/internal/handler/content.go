package handler

import (
	"strconv"
	"strings"

	"forwhis-infra/apps/api/internal/middleware"
	"forwhis-infra/apps/api/internal/models"
	"forwhis-infra/apps/api/internal/request"
	"forwhis-infra/apps/api/internal/service"
	rbacSvc "forwhis-infra/apps/api/internal/service/rbac"
	"forwhis-infra/apps/api/pkg/constants"
	"forwhis-infra/apps/api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

func contentLabel(contentType string) string {
	switch contentType {
	case models.TypePost:
		return "帖子"
	case models.TypeArticle:
		return "文章"
	default:
		return "页面"
	}
}

func GetContents(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	status := c.Query("status")
	categoryMid, _ := strconv.Atoi(c.Query("categoryMid"))
	contentType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	result := service.GetContentsList(page, limit, status, 0, categoryMid, contentType)
	return middleware.JSONSuccess(c, "获取成功喵", result)
}

func CreateContent(c fiber.Ctx) error {
	var req request.CreateContentRequest

	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	contentType := req.Type
	if contentType == "" {
		contentType = models.TypePost
	}
	if !models.IsContentType(contentType) {
		return middleware.JSONError(c, "内容类型无效喵", nil)
	}
	label := contentLabel(contentType)

	text := strings.TrimSpace(req.Text)
	if text == "" {
		return middleware.JSONError(c, label+"内容不能为空喵", nil)
	}
	if utils.CountRunes(text) > models.TextLimit(contentType) {
		return middleware.JSONError(c, label+"内容不能超过"+strconv.Itoa(models.TextLimit(contentType))+"字喵", nil)
	}

	title := strings.TrimSpace(req.Title)
	if contentType == models.TypePost {
		title = ""
	} else {
		if title == "" {
			return middleware.JSONError(c, label+"标题不能为空喵", nil)
		}
		if len(title) > 150 {
			return middleware.JSONError(c, "标题长度不能超过150个字符喵", nil)
		}
	}

	user := middleware.GetCurrentUserFromContext(c)
	status := req.Status
	if status == "" {
		status = "draft"
	}

	content := service.CreateContent(
		title,
		text,
		user.Uid,
		contentType,
		status,
		req.Password,
		service.ResolveCategoryMID(req.CategoryMid),
		req.TagIds,
	)

	if content == nil {
		return middleware.JSONSystemError(c, label+"创建失败喵", nil)
	}

	detail := "创建了" + label + "喵"
	if title != "" {
		detail = "创建了" + label + ": " + title + "喵"
	}
	service.CreateUserinfo(user.Uid, "createContent", detail)
	return middleware.JSONSuccess(c, label+"创建成功喵", map[string]interface{}{
		"cid":     content.Cid,
		"title":   content.Title,
		"type":    content.Type,
		"created": content.Created,
	})

}

func GetContent(c fiber.Ctx) error {
	cid, err := strconv.Atoi(c.Params("cid"))
	if err != nil {
		return middleware.JSONError(c, "文章ID无效喵", nil)
	}

	content := service.GetContentByCid(cid)
	if content == nil {
		return middleware.JSONError(c, "文章不存在喵", nil)
	}

	if content.Password != "" {
		user := middleware.GetCurrentUser(c)
		isOwner := user != nil && user.Uid == content.AuthorID
		if !isOwner {
			var req request.ContentPasswordRequest
			if err := c.Bind().JSON(&req); err != nil || req.Password != content.Password {
				return middleware.JSONError(c, "该文章需要密码访问喵", map[string]interface{}{
					"needPassword": true,
					"cid":          content.Cid,
					"title":        content.Title,
				})

			}
		}
	}

	service.IncrementViews(cid)

	category := service.GetContentCategory(cid)
	author := service.GetContentAuthor(content.AuthorID)
	result := map[string]interface{}{
		"cid":          content.Cid,
		"title":        content.Title,
		"text":         content.Text,
		"created":      content.Created,
		"modified":     content.Modified,
		"type":         content.Type,
		"status":       content.Status,
		"author_id":    content.AuthorID,
		"hasPassword":  content.Password != "",
		"views":        content.Views + 1,
		"comments_num": service.GetCommentCount(cid),
		"tags":         service.GetContentTags(cid),
	}
	if category != nil {
		result["category"] = map[string]interface{}{
			"mid":  category.Mid,
			"name": category.Name,
		}
		result["categoryMid"] = category.Mid
	}
	if author != nil {
		result["author"] = author
	}
	return middleware.JSONSuccess(c, "获取成功喵", result)
}

func UpdateContent(c fiber.Ctx) error {
	cid, err := strconv.Atoi(c.Params("cid"))
	if err != nil {
		return middleware.JSONError(c, "文章ID无效喵", nil)
	}

	var req map[string]interface{}
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	content := service.GetContentByCid(cid)
	if content == nil {
		return middleware.JSONError(c, "文章不存在喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)
	canEdit := user.Uid == content.AuthorID || rbacSvc.HasPermissionConst(user, constants.PermPostEdit)
	if !canEdit {
		return middleware.JSONError(c, "无权修改此文章喵", nil)
	}

	contentType := content.Type
	if rawType, ok := req["type"].(string); ok && rawType != "" {
		if !models.IsContentType(rawType) {
			return middleware.JSONError(c, "内容类型无效喵", nil)
		}
		contentType = rawType
	}
	label := contentLabel(contentType)

	text := content.Text
	if rawText, ok := req["text"].(string); ok {
		text = rawText
	}
	if strings.TrimSpace(text) == "" {
		return middleware.JSONError(c, label+"内容不能为空喵", nil)
	}
	if utils.CountRunes(text) > models.TextLimit(contentType) {
		return middleware.JSONError(c, label+"内容不能超过"+strconv.Itoa(models.TextLimit(contentType))+"字喵", nil)
	}

	if contentType == models.TypePost {
		req["title"] = ""
	} else {
		title := content.Title
		if rawTitle, ok := req["title"].(string); ok {
			title = rawTitle
		}
		title = strings.TrimSpace(title)
		if title == "" {
			return middleware.JSONError(c, label+"标题不能为空喵", nil)
		}
		if len(title) > 150 {
			return middleware.JSONError(c, "标题长度不能超过150个字符喵", nil)
		}
		req["title"] = title
	}

	if categoryMid, ok := req["categoryMid"].(float64); ok && int(categoryMid) == 0 {
		req["categoryMid"] = float64(service.ResolveCategoryMID(0))
	}
	req["type"] = contentType

	updated := service.UpdateContent(cid, req)
	if updated == nil {
		return middleware.JSONSystemError(c, label+"更新失败喵", nil)
	}

	detail := "更新了" + label + "喵"
	if updated.Title != "" {
		detail = "更新了" + label + ": " + updated.Title + "喵"
	}
	service.CreateUserinfo(user.Uid, "updateContent", detail)
	return middleware.JSONSuccess(c, label+"更新成功喵", map[string]interface{}{
		"cid":      updated.Cid,
		"title":    updated.Title,
		"type":     updated.Type,
		"modified": updated.Modified,
	})

}

func DeleteContent(c fiber.Ctx) error {
	cid, err := strconv.Atoi(c.Params("cid"))
	if err != nil {
		return middleware.JSONError(c, "文章ID无效喵", nil)
	}

	content := service.GetContentByCid(cid)
	if content == nil {
		return middleware.JSONError(c, "文章不存在喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)
	canDelete := user.Uid == content.AuthorID || rbacSvc.HasPermissionConst(user, constants.PermPostDelete)
	if !canDelete {
		return middleware.JSONError(c, "无权删除此文章喵", nil)
	}

	success := service.DeleteContent(cid)
	if !success {
		return middleware.JSONError(c, "文章删除失败喵", nil)
	}

	service.CreateUserinfo(user.Uid, "deleteContent", "删除了文章: "+strconv.Itoa(cid)+"喵")
	return middleware.JSONSuccess(c, "文章删除成功喵", nil)
}

func GetMyContents(c fiber.Ctx) error {
	user := middleware.GetCurrentUserFromContext(c)
	contents := service.GetContentsByAuthor(user.Uid)
	return middleware.JSONSuccess(c, "获取成功喵", contents)
}
