package controller

import (
	"strconv"
	"strings"

	"forestwhisperGo/magic/constants"
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/service"
	rbacSvc "forestwhisperGo/magic/service/rbac"

	"github.com/gin-gonic/gin"
)

func GetContents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	result := service.GetContentsList(page, limit, status, 0)
	middleware.JSONSuccess(c, "获取成功喵", result)
}

func CreateContent(c *gin.Context) {
	var req struct {
		Title       string `json:"title"`
		Text        string `json:"text"`
		Type        string `json:"type"`
		Status      string `json:"status"`
		Password    string `json:"password"`
		Template    string `json:"template"`
		CategoryMid int    `json:"categoryMid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		middleware.JSONError(c, "文章标题不能为空喵", nil)
		return
	}
	if strings.TrimSpace(req.Text) == "" {
		middleware.JSONError(c, "文章内容不能为空喵", nil)
		return
	}
	if req.CategoryMid == 0 {
		middleware.JSONError(c, "请选择分类喵", nil)
		return
	}
	if len(req.Title) > 150 {
		middleware.JSONError(c, "标题长度不能超过150个字符喵", nil)
		return
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
		req.Template,
		req.CategoryMid,
	)

	if content == nil {
		middleware.JSONSystemError(c, "文章创建失败喵", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "createContent", "创建了文章: "+req.Title+"喵")
	middleware.JSONSuccess(c, "文章创建成功喵", map[string]interface{}{
		"cid":     content.Cid,
		"title":   content.Title,
		"created": content.Created,
	})
}

func GetContent(c *gin.Context) {
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		middleware.JSONError(c, "文章ID无效喵", nil)
		return
	}

	content := service.GetContentByCid(cid)
	if content == nil {
		middleware.JSONError(c, "文章不存在喵", nil)
		return
	}

	if content.Password != "" {
		user := middleware.GetCurrentUser(c)
		isOwner := user != nil && user.Uid == content.AuthorID
		if !isOwner {
			var req struct {
				Password string `json:"password"`
			}
			if err := c.ShouldBindJSON(&req); err != nil || req.Password != content.Password {
				middleware.JSONError(c, "该文章需要密码访问喵", map[string]interface{}{
					"needPassword": true,
					"cid":          content.Cid,
					"title":        content.Title,
				})
				return
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
		"password":     content.Password,
		"views":        content.Views + 1,
		"comments_num": content.CommentsNum,
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
	middleware.JSONSuccess(c, "获取成功喵", result)
}

func UpdateContent(c *gin.Context) {
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		middleware.JSONError(c, "文章ID无效喵", nil)
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	content := service.GetContentByCid(cid)
	if content == nil {
		middleware.JSONError(c, "文章不存在喵", nil)
		return
	}

	user := middleware.GetCurrentUserFromContext(c)
	canEdit := user.Uid == content.AuthorID || rbacSvc.HasPermissionConst(user, constants.PermPostEdit)
	if !canEdit {
		middleware.JSONError(c, "无权修改此文章喵", nil)
		return
	}

	if title, ok := req["title"].(string); ok {
		if strings.TrimSpace(title) == "" {
			middleware.JSONError(c, "文章标题不能为空喵", nil)
			return
		}
		if len(title) > 150 {
			middleware.JSONError(c, "标题长度不能超过150个字符喵", nil)
			return
		}
	}

	updated := service.UpdateContent(cid, req)
	if updated == nil {
		middleware.JSONSystemError(c, "文章更新失败喵", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "updateContent", "更新了文章: "+updated.Title+"喵")
	middleware.JSONSuccess(c, "文章更新成功喵", map[string]interface{}{
		"cid":      updated.Cid,
		"title":    updated.Title,
		"modified": updated.Modified,
	})
}

func DeleteContent(c *gin.Context) {
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		middleware.JSONError(c, "文章ID无效喵", nil)
		return
	}

	content := service.GetContentByCid(cid)
	if content == nil {
		middleware.JSONError(c, "文章不存在喵", nil)
		return
	}

	user := middleware.GetCurrentUserFromContext(c)
	canDelete := user.Uid == content.AuthorID || rbacSvc.HasPermissionConst(user, constants.PermPostDelete)
	if !canDelete {
		middleware.JSONError(c, "无权删除此文章喵", nil)
		return
	}

	success := service.DeleteContent(cid)
	if !success {
		middleware.JSONError(c, "文章删除失败喵", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "deleteContent", "删除了文章: "+strconv.Itoa(cid)+"喵")
	middleware.JSONSuccess(c, "文章删除成功喵", nil)
}

func GetMyContents(c *gin.Context) {
	user := middleware.GetCurrentUserFromContext(c)
	contents := service.GetContentsByAuthor(user.Uid)
	middleware.JSONSuccess(c, "获取成功喵", contents)
}
