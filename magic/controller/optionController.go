package controller

import (
	"strings"

	"forestwhisperGo/magic/constants"
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/service"
	rbacSvc "forestwhisperGo/magic/service/rbac"

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

	middleware.JSONSuccess(c, "获取成功", options)
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

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
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

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空", nil)
		return
	}

	if len(req) == 0 {
		middleware.JSONError(c, "请求数据不能为空", nil)
		return
	}

	for key := range req {
		if strings.TrimSpace(key) == "" {
			middleware.JSONError(c, "选项键不能为空", nil)
			return
		}
		if len(key) > 150 {
			middleware.JSONError(c, "选项键长度不能超过150个字符", nil)
			return
		}
	}

	if err := service.SetOptions(req); err != nil {
		middleware.JSONSystemError(c, "设置保存失败", nil)
		return
	}

	service.CreateUserinfo(user.Uid, "setOptions", "批量修改了系统设置")

	middleware.JSONSuccess(c, "设置保存成功", req)
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
