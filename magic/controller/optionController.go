package controller

import (
	"strings"

	"forestwhisperGo/magic/interfaces/http/middleware"
	"forestwhisperGo/magic/models/types"
	"forestwhisperGo/magic/service"
	rbacSvc "forestwhisperGo/magic/service/rbac"
	"forestwhisperGo/pkg/constants"

	"github.com/gin-gonic/gin"
)

func GetOptions(c *gin.Context) {
	prefix := c.Query("prefix")

	var options map[string]string
	if prefix != "" {
		options = service.GetOptionsByPrefix(prefix)
	} else {
		options = service.GetAllOptions()
	}

	flatOptions := flattenOptions(options)

	middleware.JSONSuccess(c, "获取成功", flatOptions)
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

func GetOption(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		middleware.JSONError(c, "选项键不能为空", nil)
		return
	}

	value := service.GetOption(key, "")
	if value == "" {
		middleware.JSONError(c, "选项不存在", nil)
		return
	}

	middleware.JSONSuccess(c, "获取成功", map[string]interface{}{
		"key":   key,
		"value": value,
	})
}

func SetOption(c *gin.Context) {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermSettingsEdit) {
		middleware.JSONError(c, "无权修改系统设置", nil)
		return
	}

	var req types.SetOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空", nil)
		return
	}

	if strings.TrimSpace(req.Key) == "" {
		middleware.JSONError(c, "选项键不能为空", nil)
		return
	}

	if len(req.Key) > 150 {
		middleware.JSONError(c, "选项键长度不能超过150个字符", nil)
		return
	}

	if err := service.SetOption(req.Key, req.Value); err != nil {
		middleware.JSONSystemError(c, "设置保存失败", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "setOption", "修改了系统设置: "+req.Key)

	middleware.JSONSuccess(c, "设置保存成功", map[string]interface{}{
		"key":   req.Key,
		"value": req.Value,
	})
}

func SetOptions(c *gin.Context) {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermSettingsEdit) {
		middleware.JSONError(c, "无权修改系统设置", nil)
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空", nil)
		return
	}

	if len(req) == 0 {
		middleware.JSONError(c, "请求数据不能为空", nil)
		return
	}

	flatOptions := expandOptions(req)

	for key := range flatOptions {
		if strings.TrimSpace(key) == "" {
			middleware.JSONError(c, "选项键不能为空", nil)
			return
		}
		if len(key) > 150 {
			middleware.JSONError(c, "选项键长度不能超过150个字符", nil)
			return
		}
	}

	if err := service.SetOptions(flatOptions); err != nil {
		middleware.JSONSystemError(c, "设置保存失败", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "setOptions", "批量修改了系统设置")

	middleware.JSONSuccess(c, "设置保存成功", req)
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

func DeleteOption(c *gin.Context) {
	user := middleware.GetCurrentUserFromContext(c)
	if !rbacSvc.HasPermissionConst(user, constants.PermSettingsEdit) {
		middleware.JSONError(c, "无权删除系统设置", nil)
		return
	}

	key := c.Param("key")
	if key == "" {
		middleware.JSONError(c, "选项键不能为空", nil)
		return
	}

	if !service.DeleteOption(key) {
		middleware.JSONError(c, "选项不存在或删除失败", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "deleteOption", "删除了系统设置: "+key)

	middleware.JSONSuccess(c, "选项删除成功", nil)
}
