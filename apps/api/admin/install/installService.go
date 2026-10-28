package install

import (
	"fmt"

	"forestwhisperGo/apps/api/internal/models"
	"forestwhisperGo/apps/api/internal/service"
	"forestwhisperGo/apps/api/pkg/db"
	"forestwhisperGo/apps/api/pkg/utils"
	"forestwhisperGo/apps/api/pkg/utils/configutil"

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
		DbUser: configutil.GetString(dbConfig, "dbUser"),
		DbPass: configutil.GetString(dbConfig, "dbPassword"),
		DbHost: configutil.GetString(dbConfig, "dbHost"),
		DbPort: configutil.GetInt(dbConfig, "dbPort"),
		DbName: configutil.GetString(dbConfig, "dbName"),
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
	utils.SetConfig("db", "PGSQLHOST", configutil.GetString(dbConfig, "dbHost"))
	utils.SetConfig("db", "PGSQLPORT", configutil.GetInt(dbConfig, "dbPort"))
	utils.SetConfig("db", "PGSQLUSER", configutil.GetString(dbConfig, "dbUser"))
	utils.SetConfig("db", "PGSQLPWD", configutil.GetString(dbConfig, "dbPassword"))
	utils.SetConfig("db", "PGSQL_DB", configutil.GetString(dbConfig, "dbName"))
	utils.SetConfig("jwt", "JWT_ISS", "lmoadll")
	utils.SetConfig("jwt", "JWT_AUD", "lmoadll")
	return nil
}

func InitDatabase() error {
	if err := db.ReloadDB(); err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("数据库未连接")
	}
	return database.AutoMigrate(models.All()...)
}

func CreateAdminUser(username, email, password string) error {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}
	user, err := service.CreateUser(username, email, passwordHash, "", "superadmin", false)
	if err != nil {
		return fmt.Errorf("创建管理员失败: %v", err)
	}
	if user == nil {
		return fmt.Errorf("创建管理员失败")
	}
	return nil
}
