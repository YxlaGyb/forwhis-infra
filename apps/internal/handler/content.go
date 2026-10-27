package handler

import (
	"strconv"
	"strings"

	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/request"
	"forestwhisperGo/apps/internal/service"
	rbacSvc "forestwhisperGo/apps/internal/service/rbac"
	"forestwhisperGo/apps/pkg/constants"

	"github.com/gofiber/fiber/v3"
)

func GetContents(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	status := c.Query("status")
	categoryMid, _ := strconv.Atoi(c.Query("categoryMid"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	result := service.GetContentsList(page, limit, status, 0, categoryMid)
	return middleware.JSONSuccess(c, "获取成功喵", result)
}

func CreateContent(c fiber.Ctx) error {
	var req request.CreateContentRequest

	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}
	if strings.TrimSpace(req.Title) == "" {
		return middleware.JSONError(c, "文章标题不能为空喵", nil)
	}
	if strings.TrimSpace(req.Text) == "" {
		return middleware.JSONError(c, "文章内容不能为空喵", nil)
	}
	if req.CategoryMid == 0 {
		return middleware.JSONError(c, "请选择分类喵", nil)
	}
	if len(req.Title) > 150 {
		return middleware.JSONError(c, "标题长度不能超过150个字符喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)
	contentType := req.Type
	if contentType == "" {
		contentType = "post"
	}
	status := req.Status
	if status == "" {
		status = "draft"
	}

	content := service.CreateContent(
		strings.TrimSpace(req.Title),
		req.Text,
		user.Uid,
		contentType,
		status,
		req.Password,
		req.CategoryMid,
		req.TagIds,
	)

	if content == nil {
		return middleware.JSONSystemError(c, "文章创建失败喵", nil)
	}

	service.CreateUserinfo(user.Uid, "createContent", "创建了文章: "+req.Title+"喵")
	return middleware.JSONSuccess(c, "文章创建成功喵", map[string]interface{}{
		"cid":     content.Cid,
		"title":   content.Title,
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

	if title, ok := req["title"].(string); ok {
		if strings.TrimSpace(title) == "" {
			return middleware.JSONError(c, "文章标题不能为空喵", nil)
		}
		if len(title) > 150 {
			return middleware.JSONError(c, "标题长度不能超过150个字符喵", nil)
		}
	}

	updated := service.UpdateContent(cid, req)
	if updated == nil {
		return middleware.JSONSystemError(c, "文章更新失败喵", nil)
	}

	service.CreateUserinfo(user.Uid, "updateContent", "更新了文章: "+updated.Title+"喵")
	return middleware.JSONSuccess(c, "文章更新成功喵", map[string]interface{}{
		"cid":      updated.Cid,
		"title":    updated.Title,
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
