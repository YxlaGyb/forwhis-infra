package controller

import (
	"strconv"
	"strings"

	"forestwhisperGo/magic/interfaces/http/middleware"
	"forestwhisperGo/magic/models/types"
	"forestwhisperGo/magic/service"
	rbacSvc "forestwhisperGo/magic/service/rbac"
	"forestwhisperGo/pkg/constants"

	"github.com/gin-gonic/gin"
)

func GetComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
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
	middleware.JSONSuccess(c, "获取成功喵", result)
}

func GetMyComments(c *gin.Context) {
	user := middleware.GetCurrentUserFromContext(c)
	comments := service.GetCommentsByAuthor(user.Uid)
	middleware.JSONSuccess(c, "获取成功喵", comments)
}

func CreateComment(c *gin.Context) {
	var req types.CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}
	if req.Cid == 0 {
		middleware.JSONError(c, "文章ID不能为空喵", nil)
		return
	}
	if strings.TrimSpace(req.Text) == "" {
		middleware.JSONError(c, "评论内容不能为空喵", nil)
		return
	}
	if len(req.Text) > 5000 {
		middleware.JSONError(c, "评论长度不能超过5000个字符喵", nil)
		return
	}

	content := service.GetContentByCid(req.Cid)
	if content == nil {
		middleware.JSONError(c, "文章不存在喵", nil)
		return
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

	middleware.JSONSuccess(c, "评论已提交, 等待审核喵", map[string]interface{}{
		"coid":    comment.Coid,
		"created": comment.Created,
	})
}

func ApproveComment(c *gin.Context) {
	coid, err := strconv.Atoi(c.Param("coid"))
	if err != nil {
		middleware.JSONError(c, "评论ID无效喵", nil)
		return
	}

	comment := service.GetCommentByCoid(coid)
	if comment == nil {
		middleware.JSONError(c, "评论不存在喵", nil)
		return
	}

	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermCommentModerate) {
		middleware.JSONError(c, "无权审核此评论喵", nil)
		return
	}

	updated := service.ApproveComment(coid)
	if updated == nil {
		middleware.JSONError(c, "评论审核失败喵", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "approveComment", "审核通过了评论: "+strconv.Itoa(coid)+"喵")
	middleware.JSONSuccess(c, "评论审核通过喵", nil)
}

func DeleteComment(c *gin.Context) {
	coid, err := strconv.Atoi(c.Param("coid"))
	if err != nil {
		middleware.JSONError(c, "评论ID无效喵", nil)
		return
	}

	comment := service.GetCommentByCoid(coid)
	if comment == nil {
		middleware.JSONError(c, "评论不存在喵", nil)
		return
	}

	user := middleware.GetCurrentUserFromContext(c)
	canDelete := comment.AuthorID == user.Uid || rbacSvc.HasPermissionConst(user, constants.PermCommentDelete)

	if !canDelete {
		middleware.JSONError(c, "无权删除此评论喵", nil)
		return
	}

	success := service.DeleteComment(coid)
	if !success {
		middleware.JSONError(c, "评论删除失败喵", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "deleteComment", "删除了评论: "+strconv.Itoa(coid)+"喵")
	middleware.JSONSuccess(c, "评论删除成功喵", nil)
}

func GetComment(c *gin.Context) {
	coid, err := strconv.Atoi(c.Param("coid"))
	if err != nil {
		middleware.JSONError(c, "评论ID无效喵", nil)
		return
	}

	comment := service.GetCommentByCoid(coid)
	if comment == nil {
		middleware.JSONError(c, "评论不存在喵", nil)
		return
	}

	middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
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

func UpdateComment(c *gin.Context) {
	coid, err := strconv.Atoi(c.Param("coid"))
	if err != nil {
		middleware.JSONError(c, "评论ID无效喵", nil)
		return
	}

	comment := service.GetCommentByCoid(coid)
	if comment == nil {
		middleware.JSONError(c, "评论不存在喵", nil)
		return
	}

	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermCommentModerate) {
		middleware.JSONError(c, "无权修改此评论喵", nil)
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	updated := service.UpdateComment(coid, req)
	if updated == nil {
		middleware.JSONError(c, "评论修改失败喵", nil)
		return
	}

	middleware.JSONSuccess(c, "评论修改成功喵", map[string]interface{}{
		"coid":   updated.Coid,
		"status": updated.Status,
	})
}

func GetContentComments(c *gin.Context) {
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
	middleware.JSONSuccess(c, "获取成功喵", comments)
}
