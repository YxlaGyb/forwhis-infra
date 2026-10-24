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

func GetMetas(c fiber.Ctx) error {
	metaType := c.Query("type")

	if metaType != "" && metaType != "category" && metaType != "tag" {
		return middleware.JSONError(c, "类型参数无效喵", nil)
	}

	metas := service.GetMetasList(metaType)
	return middleware.JSONSuccess(c, "获取成功喵", metas)
}

func CreateMeta(c fiber.Ctx) error {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermPostCreate) {
		return middleware.JSONError(c, "无权创建分类喵", nil)
	}

	var req request.CreateMetaRequest
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	if strings.TrimSpace(req.Name) == "" {
		return middleware.JSONError(c, "名称不能为空喵", nil)
	}
	if len(req.Name) > 150 {
		return middleware.JSONError(c, "名称长度不能超过150个字符喵", nil)
	}

	metaType := req.Type
	if metaType == "" {
		metaType = "category"
	}
	if metaType != "category" && metaType != "tag" {
		return middleware.JSONError(c, "类型必须为 category 或 tag 喵", nil)
	}

	meta := service.CreateMeta(
		strings.TrimSpace(req.Name),
		metaType,
		req.Slug,
		req.Description,
		req.Parent,
	)

	service.CreateUserinfo(user.Uid, "createMeta", "创建了分类: "+req.Name+"喵")
	return middleware.JSONSuccess(c, "分类创建成功喵", map[string]interface{}{
		"mid":  meta.Mid,
		"name": meta.Name,
		"slug": meta.Slug,
		"type": meta.Type,
	})

}

func UpdateMeta(c fiber.Ctx) error {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermPostEdit) {
		return middleware.JSONError(c, "无权修改分类喵", nil)
	}

	mid, err := strconv.Atoi(c.Params("mid"))
	if err != nil {
		return middleware.JSONError(c, "分类ID无效喵", nil)
	}

	var req map[string]interface{}
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	meta := service.GetMetaByMid(mid)
	if meta == nil {
		return middleware.JSONError(c, "分类不存在喵", nil)
	}

	if name, ok := req["name"].(string); ok {
		if strings.TrimSpace(name) == "" {
			return middleware.JSONError(c, "名称不能为空喵", nil)
		}
		if len(name) > 150 {
			return middleware.JSONError(c, "名称长度不能超过150个字符喵", nil)
		}
	}

	updated := service.UpdateMeta(mid, req)
	if updated == nil {
		return middleware.JSONError(c, "分类更新失败喵", nil)
	}

	service.CreateUserinfo(user.Uid, "updateMeta", "更新了分类: "+updated.Name+"喵")
	return middleware.JSONSuccess(c, "分类更新成功喵", map[string]interface{}{
		"mid":  updated.Mid,
		"name": updated.Name,
	})

}

func DeleteMeta(c fiber.Ctx) error {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermPostDelete) {
		return middleware.JSONError(c, "无权删除分类喵", nil)
	}

	mid, err := strconv.Atoi(c.Params("mid"))
	if err != nil {
		return middleware.JSONError(c, "分类ID无效喵", nil)
	}

	meta := service.GetMetaByMid(mid)
	if meta == nil {
		return middleware.JSONError(c, "分类不存在喵", nil)
	}

	success := service.DeleteMeta(mid)
	if !success {
		return middleware.JSONSystemError(c, "分类删除失败喵", nil)
	}

	service.CreateUserinfo(user.Uid, "deleteMeta", "删除了分类: "+strconv.Itoa(mid)+"喵")
	return middleware.JSONSuccess(c, "分类删除成功喵", nil)
}

func GetMeta(c fiber.Ctx) error {
	mid, err := strconv.Atoi(c.Params("mid"))
	if err != nil {
		return middleware.JSONError(c, "分类ID无效喵", nil)
	}

	meta := service.GetMetaByMid(mid)
	if meta == nil {
		return middleware.JSONError(c, "分类不存在喵", nil)
	}
	return middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
		"mid":         meta.Mid,
		"name":        meta.Name,
		"slug":        meta.Slug,
		"type":        meta.Type,
		"description": meta.Description,
		"count":       meta.Count,
		"parent":      meta.Parent,
	})

}
