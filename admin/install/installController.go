package install

import (
	"net/http"
	"os"

	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/utils"

	"github.com/gin-gonic/gin"
)

func GetInstallHtml(c *gin.Context) {
	htmlPath := "admin/install/install.html"
	content, err := os.ReadFile(htmlPath)
	if err != nil {
		c.Data(http.StatusNotFound, "text/plain", []byte("安装页面不存在"))
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}

func InstallIndex(c *gin.Context) {
	GetInstallHtml(c)
}

func InstallPage(c *gin.Context) {
	GetInstallHtml(c)
}

func CheckInstallStatus(c *gin.Context) {
	middleware.JSONSuccess(c, "获取成功", gin.H{
		"installed":    utils.IsInstalled(),
		"configExists": utils.ConfigExists(),
	})
}

func TestDbConnectionHandler(c *gin.Context) {
	var req struct {
		DbType   string                 `json:"dbType"`
		DbConfig map[string]interface{} `json:"dbConfig"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	success, msg := TestDbConnection(req.DbType, req.DbConfig)
	if success {
		middleware.JSONSuccess(c, msg, nil)
	} else {
		middleware.JSONError(c, msg, nil)
	}
}

func InstallSite(c *gin.Context) {
	var req struct {
		SiteName      string                 `json:"siteName"`
		SiteUrl       string                 `json:"siteUrl"`
		AdminUsername string                 `json:"adminUsername"`
		AdminEmail    string                 `json:"adminEmail"`
		AdminPassword string                 `json:"adminPassword"`
		DbType        string                 `json:"dbType"`
		DbConfig      map[string]interface{} `json:"dbConfig"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	if req.SiteName == "" || req.AdminUsername == "" || req.AdminEmail == "" || req.AdminPassword == "" {
		middleware.JSONError(c, "请填写完整信息喵", nil)
		return
	}

	if !utils.IsValidName(req.AdminUsername) {
		middleware.JSONError(c, "用户名格式不正确喵", nil)
		return
	}

	if !utils.IsValidEmail(req.AdminEmail) {
		middleware.JSONError(c, "邮箱格式不正确喵", nil)
		return
	}

	if !utils.IsValidPassword(req.AdminPassword) {
		middleware.JSONError(c, "密码长度应不少于8个字符喵", nil)
		return
	}

	if req.DbType == "postgresql" {
		success, msg := TestDbConnection(req.DbType, req.DbConfig)
		if !success {
			middleware.JSONError(c, "数据库连接失败: "+msg, nil)
			return
		}
	}

	if err := SaveConfig(req.SiteName, req.SiteUrl, req.DbType, req.DbConfig); err != nil {
		middleware.JSONError(c, "配置保存失败: "+err.Error(), nil)
		return
	}

	if err := InitDatabase(); err != nil {
		middleware.JSONError(c, "数据库初始化失败: "+err.Error(), nil)
		return
	}

	if err := CreateAdminUser(req.AdminUsername, req.AdminEmail, req.AdminPassword); err != nil {
		middleware.JSONError(c, "创建管理员失败: "+err.Error(), nil)
		return
	}

	middleware.JSONSuccess(c, "安装成功喵", nil)
}
