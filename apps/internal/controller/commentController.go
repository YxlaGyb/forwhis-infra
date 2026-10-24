package controller

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

func GetComments(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	cid, _ := strconv.Atoi(c.Query("cid"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	result := service.GetCommentsList(cid, status, 0, page, limit)
	return middleware.JSONSuccess(c, "获取成功喵", result)
}

func GetMyComments(c fiber.Ctx) error {
	user := middleware.GetCurrentUserFromContext(c)
	comments := service.GetCommentsByAuthor(user.Uid)
	return middleware.JSONSuccess(c, "获取成功喵", comments)
}

func CreateComment(c fiber.Ctx) error {
	var req request.CreateCommentRequest

	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}
	if req.Cid == 0 {
		return middleware.JSONError(c, "文章ID不能为空喵", nil)
	}
	if strings.TrimSpace(req.Text) == "" {
		return middleware.JSONError(c, "评论内容不能为空喵", nil)
	}
	if len(req.Text) > 5000 {
		return middleware.JSONError(c, "评论长度不能超过5000个字符喵", nil)
	}

	content := service.GetContentByCid(req.Cid)
	if content == nil {
		return middleware.JSONError(c, "文章不存在喵", nil)
	}

	user := middleware.GetCurrentUser(c)
	var authorID int
	var author, mail string
	if user != nil {
		authorID = user.Uid
		author = user.Name
		mail = user.Mail
	}

	commentType := req.Type
	if commentType == "" {
		commentType = "comment"
	}

	comment := service.CreateComment(req.Cid, strings.TrimSpace(req.Text), authorID, author, mail, req.Url, commentType, req.Parent)

	if authorID > 0 {
		service.CreateUserinfo(authorID, "createComment", "发表了评论喵")
	}
	return middleware.JSONSuccess(c, "评论已提交, 等待审核喵", map[string]interface{}{
		"coid":    comment.Coid,
		"created": comment.Created,
	})

}

func ApproveComment(c fiber.Ctx) error {
	coid, err := strconv.Atoi(c.Params("coid"))
	if err != nil {
		return middleware.JSONError(c, "评论ID无效喵", nil)
	}

	comment := service.GetCommentByCoid(coid)
	if comment == nil {
		return middleware.JSONError(c, "评论不存在喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermCommentModerate) {
		return middleware.JSONError(c, "无权审核此评论喵", nil)
	}

	updated := service.ApproveComment(coid)
	if updated == nil {
		return middleware.JSONError(c, "评论审核失败喵", nil)
	}

	service.CreateUserinfo(user.Uid, "approveComment", "审核通过了评论: "+strconv.Itoa(coid)+"喵")
	return middleware.JSONSuccess(c, "评论审核通过喵", nil)
}

func DeleteComment(c fiber.Ctx) error {
	coid, err := strconv.Atoi(c.Params("coid"))
	if err != nil {
		return middleware.JSONError(c, "评论ID无效喵", nil)
	}

	comment := service.GetCommentByCoid(coid)
	if comment == nil {
		return middleware.JSONError(c, "评论不存在喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)
	canDelete := comment.AuthorID == user.Uid || rbacSvc.HasPermissionConst(user, constants.PermCommentDelete)

	if !canDelete {
		return middleware.JSONError(c, "无权删除此评论喵", nil)
	}

	success := service.DeleteComment(coid)
	if !success {
		return middleware.JSONError(c, "评论删除失败喵", nil)
	}

	service.CreateUserinfo(user.Uid, "deleteComment", "删除了评论: "+strconv.Itoa(coid)+"喵")
	return middleware.JSONSuccess(c, "评论删除成功喵", nil)
}

func GetComment(c fiber.Ctx) error {
	coid, err := strconv.Atoi(c.Params("coid"))
	if err != nil {
		return middleware.JSONError(c, "评论ID无效喵", nil)
	}

	comment := service.GetCommentByCoid(coid)
	if comment == nil {
		return middleware.JSONError(c, "评论不存在喵", nil)
	}
	return middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
		"coid":    comment.Coid,
		"cid":     comment.Cid,
		"created": comment.Created,
		"author":  comment.Author,
		"text":    comment.Text,
		"type":    comment.Type,
		"status":  comment.Status,
		"parent":  comment.Parent,
	})

}

func UpdateComment(c fiber.Ctx) error {
	coid, err := strconv.Atoi(c.Params("coid"))
	if err != nil {
		return middleware.JSONError(c, "评论ID无效喵", nil)
	}

	comment := service.GetCommentByCoid(coid)
	if comment == nil {
		return middleware.JSONError(c, "评论不存在喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermCommentModerate) {
		return middleware.JSONError(c, "无权修改此评论喵", nil)
	}

	var req map[string]interface{}
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	updated := service.UpdateComment(coid, req)
	if updated == nil {
		return middleware.JSONError(c, "评论修改失败喵", nil)
	}
	return middleware.JSONSuccess(c, "评论修改成功喵", map[string]interface{}{
		"coid":   updated.Coid,
		"status": updated.Status,
	})

}

func GetContentComments(c fiber.Ctx) error {
	cid, err := strconv.Atoi(c.Params("cid"))
	if err != nil {
		return middleware.JSONError(c, "文章ID无效喵", nil)
	}

	content := service.GetContentByCid(cid)
	if content == nil {
		return middleware.JSONError(c, "文章不存在喵", nil)
	}

	user := middleware.GetCurrentUser(c)
	var viewerID int
	var canModerate bool
	if user != nil {
		viewerID = user.Uid
		canModerate = rbacSvc.HasPermissionConst(user, constants.PermCommentModerate)
	}

	var comments []map[string]interface{}
	if canModerate {
		comments = service.GetAllCommentsByContent(cid)
	} else {
		comments = service.GetCommentsByContent(cid, "", viewerID)
	}
	return middleware.JSONSuccess(c, "获取成功喵", comments)
}
