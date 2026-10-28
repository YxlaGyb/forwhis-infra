package handler

import (
	"strconv"
	"strings"

	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/request"
	"forestwhisperGo/apps/api/internal/service"
	rbacSvc "forestwhisperGo/apps/api/internal/service/rbac"
	"forestwhisperGo/apps/api/pkg/constants"

	"github.com/gofiber/fiber/v3"
)

func metaTypeFromQuery(c fiber.Ctx) string {
	if c.Query("type") == "tag" {
		return "tag"
	}
	return "category"
}

func GetMetas(c fiber.Ctx) error {
	metaType := c.Query("type")
	if metaType == "tag" {
		return middleware.JSONSuccess(c, "获取成功喵", service.GetTagsList())
	}
	if metaType != "" && metaType != "category" {
		return middleware.JSONError(c, "类型参数无效喵", nil)
	}
	return middleware.JSONSuccess(c, "获取成功喵", service.GetCategoriesList())
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

	if metaType == "tag" {
		tag := service.CreateTag(strings.TrimSpace(req.Name), req.Slug)
		service.CreateUserinfo(user.Uid, "createMeta", "创建了标签: "+req.Name+"喵")
		return middleware.JSONSuccess(c, "标签创建成功喵", map[string]interface{}{
			"mid": tag.Id, "id": tag.Id, "name": tag.Name, "slug": tag.Slug, "type": "tag",
		})
	}

	category := service.CreateCategory(strings.TrimSpace(req.Name), req.Slug, req.Description, req.Parent)
	service.CreateUserinfo(user.Uid, "createMeta", "创建了分类: "+req.Name+"喵")
	return middleware.JSONSuccess(c, "分类创建成功喵", map[string]interface{}{
		"mid": category.Mid, "name": category.Name, "slug": category.Slug, "type": "category",
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

	if name, ok := req["name"].(string); ok {
		if strings.TrimSpace(name) == "" {
			return middleware.JSONError(c, "名称不能为空喵", nil)
		}
		if len(name) > 150 {
			return middleware.JSONError(c, "名称长度不能超过150个字符喵", nil)
		}
	}

	metaType := metaTypeFromQuery(c)
	if t, ok := req["type"].(string); ok && t == "tag" {
		metaType = "tag"
	}

	if metaType == "tag" {
		tag := service.UpdateTag(mid, req)
		if tag == nil {
			return middleware.JSONError(c, "标签不存在喵", nil)
		}
		service.CreateUserinfo(user.Uid, "updateMeta", "更新了标签: "+tag.Name+"喵")
		return middleware.JSONSuccess(c, "标签更新成功喵", map[string]interface{}{
			"mid": tag.Id, "id": tag.Id, "name": tag.Name, "type": "tag",
		})
	}

	category := service.UpdateCategory(mid, req)
	if category == nil {
		return middleware.JSONError(c, "分类不存在喵", nil)
	}
	service.CreateUserinfo(user.Uid, "updateMeta", "更新了分类: "+category.Name+"喵")
	return middleware.JSONSuccess(c, "分类更新成功喵", map[string]interface{}{
		"mid": category.Mid, "name": category.Name, "type": "category",
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

	if metaTypeFromQuery(c) == "tag" {
		if !service.DeleteTag(mid) {
			return middleware.JSONError(c, "标签不存在喵", nil)
		}
		service.CreateUserinfo(user.Uid, "deleteMeta", "删除了标签: "+strconv.Itoa(mid)+"喵")
		return middleware.JSONSuccess(c, "标签删除成功喵", nil)
	}

	if !service.DeleteCategory(mid) {
		return middleware.JSONError(c, "分类不存在喵", nil)
	}
	service.CreateUserinfo(user.Uid, "deleteMeta", "删除了分类: "+strconv.Itoa(mid)+"喵")
	return middleware.JSONSuccess(c, "分类删除成功喵", nil)
}

func GetMeta(c fiber.Ctx) error {
	mid, err := strconv.Atoi(c.Params("mid"))
	if err != nil {
		return middleware.JSONError(c, "分类ID无效喵", nil)
	}

	if metaTypeFromQuery(c) == "tag" {
		tag := service.GetTagById(mid)
		if tag == nil {
			return middleware.JSONError(c, "标签不存在喵", nil)
		}
		return middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
			"mid": tag.Id, "id": tag.Id, "name": tag.Name, "slug": tag.Slug, "type": "tag",
		})
	}

	category := service.GetCategoryByMid(mid)
	if category == nil {
		return middleware.JSONError(c, "分类不存在喵", nil)
	}
	return middleware.JSONSuccess(c, "获取成功喵", map[string]interface{}{
		"mid":         category.Mid,
		"name":        category.Name,
		"slug":        category.Slug,
		"description": category.Description,
		"count":       service.GetCategoryCount(mid),
		"parent":      category.Parent,
		"type":        "category",
	})
}
