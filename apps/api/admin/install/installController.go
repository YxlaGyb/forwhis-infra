package install

import (
	_ "embed"
	"net/http"

	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

//go:embed install.html
var installHtml []byte

func GetInstallHtml(c fiber.Ctx) error {
	c.Status(http.StatusOK).Type("html", "utf-8")
	return c.Send(installHtml)
}

func InstallPage(c fiber.Ctx) error {
	return GetInstallHtml(c)
}

func CheckInstallStatus(c fiber.Ctx) error {
	return middleware.JSONSuccess(c, "获取成功", fiber.Map{
		"installed":    utils.IsInstalled(),
		"configExists": utils.ConfigExists(),
	})

}

func TestDbConnectionHandler(c fiber.Ctx) error {
	var req struct {
		DbType   string                 `json:"dbType"`
		DbConfig map[string]interface{} `json:"dbConfig"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	success, msg := TestDbConnection(req.DbType, req.DbConfig)
	if success {
		return middleware.JSONSuccess(c, msg, nil)
	}
	return middleware.JSONError(c, msg, nil)
}

func InstallSite(c fiber.Ctx) error {
	var req struct {
		SiteName      string                 `json:"siteName"`
		SiteUrl       string                 `json:"siteUrl"`
		AdminUsername string                 `json:"adminUsername"`
		AdminEmail    string                 `json:"adminEmail"`
		AdminPassword string                 `json:"adminPassword"`
		DbType        string                 `json:"dbType"`
		DbConfig      map[string]interface{} `json:"dbConfig"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	if req.SiteName == "" || req.AdminUsername == "" || req.AdminEmail == "" || req.AdminPassword == "" {
		return middleware.JSONError(c, "请填写完整信息喵", nil)
	}

	if !utils.IsValidName(req.AdminUsername) {
		return middleware.JSONError(c, "用户名格式不正确喵", nil)
	}

	if !utils.IsValidEmail(req.AdminEmail) {
		return middleware.JSONError(c, "邮箱格式不正确喵", nil)
	}

	if !utils.IsValidPassword(req.AdminPassword) {
		return middleware.JSONError(c, "密码长度应不少于8个字符喵", nil)
	}

	if req.DbType == "postgresql" {
		success, msg := TestDbConnection(req.DbType, req.DbConfig)
		if !success {
			return middleware.JSONError(c, "数据库连接失败: "+msg, nil)
		}
	}

	if err := SaveConfig(req.SiteName, req.SiteUrl, req.DbType, req.DbConfig); err != nil {
		return middleware.JSONError(c, "配置保存失败: "+err.Error(), nil)
	}

	if err := InitDatabase(); err != nil {
		return middleware.JSONError(c, "数据库初始化失败: "+err.Error(), nil)
	}

	if err := CreateAdminUser(req.AdminUsername, req.AdminEmail, req.AdminPassword); err != nil {
		return middleware.JSONError(c, "创建管理员失败: "+err.Error(), nil)
	}
	return middleware.JSONSuccess(c, "安装成功喵", nil)
}
