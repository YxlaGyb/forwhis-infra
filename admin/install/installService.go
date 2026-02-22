package install

import (
	"fmt"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/rbac"
	"forestwhisperGo/magic/service"
	rbacSvc "forestwhisperGo/magic/service/rbac"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	DbUser string `json:"dbUser"`
	DbPass string `json:"dbPass"`
	DbHost string `json:"dbHost"`
	DbPort int    `json:"dbPort"`
	DbName string `json:"dbName"`
}

func TestDbConnection(dbType string, dbConfig map[string]interface{}) (bool, string) {
	if dbType != "postgresql" {
		return false, "暂只支持 postgresql"
	}

	cfg := DbConfig{
		DbUser: getString(dbConfig, "dbUser"),
		DbPass: getString(dbConfig, "dbPassword"),
		DbHost: getString(dbConfig, "dbHost"),
		DbPort: getInt(dbConfig, "dbPort"),
		DbName: getString(dbConfig, "dbName"),
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DbHost, cfg.DbUser, cfg.DbPass, cfg.DbName, cfg.DbPort)

	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return false, err.Error()
	}
	return true, "数据库连接成功"
}

func SaveConfig(siteName, siteUrl, dbType string, dbConfig map[string]interface{}) error {
	utils.SetConfig("server", "siteName", siteName)
	utils.SetConfig("server", "siteUrl", siteUrl)
	utils.SetConfig("server", "install", true)
	utils.SetConfig("db", "SQLNAME", dbType)
	utils.SetConfig("db", "PGSQLHOST", getString(dbConfig, "dbHost"))
	utils.SetConfig("db", "PGSQLPORT", getInt(dbConfig, "dbPort"))
	utils.SetConfig("db", "PGSQLUSER", getString(dbConfig, "dbUser"))
	utils.SetConfig("db", "PGSQLPWD", getString(dbConfig, "dbPassword"))
	utils.SetConfig("db", "PGSQL_DB", getString(dbConfig, "dbName"))
	utils.SetConfig("jwt", "JWT_ISS", "lmoadll")
	utils.SetConfig("jwt", "JWT_AUD", "lmoadll")
	return nil
}

func InitDatabase() error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("数据库未连接")
	}
	return database.AutoMigrate(
		&models.User{},
		&models.Usermeta{},
		&models.Userinfo{},
		&models.Content{},
		&models.Comment{},
		&models.Meta{},
		&models.Relationship{},
		&models.Review{},
		&models.ContentField{},
		&models.AIConversation{},
		&models.Like{},
		&models.Option{},
		&rbac.Role{},
		&rbac.Permission{},
		&rbac.UserRole{},
		&rbac.RolePermission{},
	)
}

func CreateAdminUser(username, email, password string) error {
	rbacSvc.InitDefaultRbac()

	passwordHash := utils.HashPassword(password)
	user, err := service.CreateUser(username, email, passwordHash, "", "superadmin", false)
	if err != nil {
		return fmt.Errorf("创建管理员失败: %v", err)
	}
	if user == nil {
		return fmt.Errorf("创建管理员失败")
	}
	return nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	switch v := m[key].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	return 0
}
