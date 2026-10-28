package handler

import (
	"strings"

	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/request"
	"forestwhisperGo/apps/api/internal/service"
	rbacSvc "forestwhisperGo/apps/api/internal/service/rbac"
	"forestwhisperGo/apps/api/pkg/constants"

	"github.com/gofiber/fiber/v3"
)

func GetOptions(c fiber.Ctx) error {
	prefix := c.Query("prefix")

	var options map[string]string
	if prefix != "" {
		options = service.GetOptionsByPrefix(prefix)
	} else {
		options = service.GetAllOptions()
	}

	flatOptions := flattenOptions(options)
	return middleware.JSONSuccess(c, "获取成功", flatOptions)
}

func flattenOptions(options map[string]string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range options {
		parts := strings.Split(key, ".")
		if len(parts) == 2 {
			section, field := parts[0], parts[1]
			if _, ok := result[section]; !ok {
				result[section] = make(map[string]interface{})
			}
			if sectionMap, ok := result[section].(map[string]interface{}); ok {
				sectionMap[field] = value
			}
		} else {
			result[key] = value
		}
	}

	return result
}

func GetOption(c fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return middleware.JSONError(c, "选项键不能为空", nil)
	}

	value := service.GetOption(key, "")
	if value == "" {
		return middleware.JSONError(c, "选项不存在", nil)
	}
	return middleware.JSONSuccess(c, "获取成功", map[string]interface{}{
		"key":   key,
		"value": value,
	})

}

func SetOption(c fiber.Ctx) error {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermSettingsEdit) {
		return middleware.JSONError(c, "无权修改系统设置", nil)
	}

	var req request.SetOptionRequest
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空", nil)
	}

	if strings.TrimSpace(req.Key) == "" {
		return middleware.JSONError(c, "选项键不能为空", nil)
	}

	if len(req.Key) > 150 {
		return middleware.JSONError(c, "选项键长度不能超过150个字符", nil)
	}

	if err := service.SetOption(req.Key, req.Value); err != nil {
		return middleware.JSONSystemError(c, "设置保存失败", nil)
	}

	service.CreateUserinfo(user.Uid, "setOption", "修改了系统设置: "+req.Key)
	return middleware.JSONSuccess(c, "设置保存成功", map[string]interface{}{
		"key":   req.Key,
		"value": req.Value,
	})

}

func SetOptions(c fiber.Ctx) error {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermSettingsEdit) {
		return middleware.JSONError(c, "无权修改系统设置", nil)
	}

	var req map[string]interface{}
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空", nil)
	}

	if len(req) == 0 {
		return middleware.JSONError(c, "请求数据不能为空", nil)
	}

	flatOptions := expandOptions(req)

	for key := range flatOptions {
		if strings.TrimSpace(key) == "" {
			return middleware.JSONError(c, "选项键不能为空", nil)
		}
		if len(key) > 150 {
			return middleware.JSONError(c, "选项键长度不能超过150个字符", nil)
		}
	}

	if err := service.SetOptions(flatOptions); err != nil {
		return middleware.JSONSystemError(c, "设置保存失败", nil)
	}

	service.CreateUserinfo(user.Uid, "setOptions", "批量修改了系统设置")
	return middleware.JSONSuccess(c, "设置保存成功", req)
}

func expandOptions(nested map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for sectionKey, sectionValue := range nested {
		if sectionMap, ok := sectionValue.(map[string]interface{}); ok {
			for fieldKey, fieldValue := range sectionMap {
				key := sectionKey + "." + fieldKey
				switch v := fieldValue.(type) {
				case bool:
					if v {
						result[key] = "1"
					} else {
						result[key] = "0"
					}
				case string:
					result[key] = v
				default:
					if v != nil {
						result[key] = "1"
					} else {
						result[key] = "0"
					}
				}
			}
		} else {
			switch v := sectionValue.(type) {
			case bool:
				if v {
					result[sectionKey] = "1"
				} else {
					result[sectionKey] = "0"
				}
			case string:
				result[sectionKey] = v
			}
		}
	}

	return result
}

func DeleteOption(c fiber.Ctx) error {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermSettingsEdit) {
		return middleware.JSONError(c, "无权删除系统设置", nil)
	}

	key := c.Params("key")
	if key == "" {
		return middleware.JSONError(c, "选项键不能为空", nil)
	}

	if !service.DeleteOption(key) {
		return middleware.JSONError(c, "选项不存在或删除失败", nil)
	}

	service.CreateUserinfo(user.Uid, "deleteOption", "删除了系统设置: "+key)
	return middleware.JSONSuccess(c, "选项删除成功", nil)
}
