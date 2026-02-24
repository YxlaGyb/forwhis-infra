package controller

import (
	"strconv"
	"strings"

	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/models/types"
	"forestwhisperGo/magic/service"
	rbacSvc "forestwhisperGo/magic/service/rbac"
	"forestwhisperGo/pkg/constants"

	"github.com/gin-gonic/gin"
)

func GetMetas(c *gin.Context) {
	metaType := c.Query("type")

	if metaType != "" && metaType != "category" && metaType != "tag" {
		middleware.JSONError(c, "类型参数无效喵", nil)
		return
	}

	metas := service.GetMetasList(metaType)
	middleware.JSONSuccess(c, "获取成功喵", metas)
}

func CreateMeta(c *gin.Context) {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermPostCreate) {
		middleware.JSONError(c, "无权创建分类喵", nil)
		return
	}

	var req types.CreateMetaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		middleware.JSONError(c, "名称不能为空喵", nil)
		return
	}
	if len(req.Name) > 150 {
		middleware.JSONError(c, "名称长度不能超过150个字符喵", nil)
		return
	}

	metaType := req.Type
	if metaType == "" {
		metaType = "category"
	}
	if metaType != "category" && metaType != "tag" {
		middleware.JSONError(c, "类型必须为 category 或 tag 喵", nil)
		return
	}

	meta := service.CreateMeta(
		strings.TrimSpace(req.Name),
		metaType,
		req.Slug,
		req.Description,
		req.Parent,
	)

	service.CreateUserinfo(user.Uid, "createMeta", "创建了分类: "+req.Name+"喵")

	middleware.JSONSuccess(c, "分类创建成功喵", map[string]interface{}{
		"mid":  meta.Mid,
		"name": meta.Name,
		"slug": meta.Slug,
		"type": meta.Type,
	})
}

func UpdateMeta(c *gin.Context) {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermPostEdit) {
		middleware.JSONError(c, "无权修改分类喵", nil)
		return
	}

	mid, err := strconv.Atoi(c.Param("mid"))
	if err != nil {
		middleware.JSONError(c, "分类ID无效喵", nil)
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	meta := service.GetMetaByMid(mid)
	if meta == nil {
		middleware.JSONError(c, "分类不存在喵", nil)
		return
	}

	if name, ok := req["name"].(string); ok {
		if strings.TrimSpace(name) == "" {
			middleware.JSONError(c, "名称不能为空喵", nil)
			return
		}
		if len(name) > 150 {
			middleware.JSONError(c, "名称长度不能超过150个字符喵", nil)
			return
		}
	}

	updated := service.UpdateMeta(mid, req)
	if updated == nil {
		middleware.JSONError(c, "分类更新失败喵", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "updateMeta", "更新了分类: "+updated.Name+"喵")

	middleware.JSONSuccess(c, "分类更新成功喵", map[string]interface{}{
		"mid":  updated.Mid,
		"name": updated.Name,
	})
}

func DeleteMeta(c *gin.Context) {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermPostDelete) {
		middleware.JSONError(c, "无权删除分类喵", nil)
		return
	}

	mid, err := strconv.Atoi(c.Param("mid"))
	if err != nil {
		middleware.JSONError(c, "分类ID无效喵", nil)
		return
	}

	meta := service.GetMetaByMid(mid)
	if meta == nil {
		middleware.JSONError(c, "分类不存在喵", nil)
		return
	}

	success := service.DeleteMeta(mid)
	if !success {
		middleware.JSONSystemError(c, "分类删除失败喵", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "deleteMeta", "删除了分类: "+strconv.Itoa(mid)+"喵")

	middleware.JSONSuccess(c, "分类删除成功喵", nil)
}

func GetMeta(c *gin.Context) {
	mid, err := strconv.Atoi(c.Param("mid"))
	if err != nil {
		middleware.JSONError(c, "分类ID无效喵", nil)
		return
	}

	meta := service.GetMetaByMid(mid)
	if meta == nil {
		middleware.JSONError(c, "分类不存在喵", nil)
		return
	}

	middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
		"mid":         meta.Mid,
		"name":        meta.Name,
		"slug":        meta.Slug,
		"type":        meta.Type,
		"description": meta.Description,
		"count":       meta.Count,
		"parent":      meta.Parent,
	})
}
