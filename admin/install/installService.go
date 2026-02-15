package install

import (
	"fmt"
	"time"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/rbac"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type InstallService struct{}

type DbConfig struct {
	DbType   string `json:"dbType"`
	DbUser   string `json:"dbUser"`
	DbPass   string `json:"dbPass"`
	DbHost   string `json:"dbHost"`
	DbPort   int    `json:"dbPort"`
	DbName   string `json:"dbName"`
}

func TestDbConnection(dbType string, dbConfig map[string]interface{}) (bool, string) {
	if dbType != "postgresql" {
		return false, "暂只支持 postgresql"
	}

	cfg := DbConfig{
		DbType: dbType,
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
		&rbac.Role{},
		&rbac.Permission{},
		&rbac.UserRole{},
		&rbac.RolePermission{},
	)
}

func CreateAdminUser(username, email, password string) error {
	passwordHash := utils.HashPassword(password)
	now := int(time.Now().Unix())
	user := models.User{
		Name: username, Mail: email, Password: passwordHash,
		CreatedAt: now, LastLogin: 0, LastIp: "", IsActive: 1,
	}
	if err := db.GetDB().Create(&user).Error; err != nil {
		return err
	}
	InitDefaultRbac(user.Uid)
	createDefaultUsermeta(user.Uid)
	return nil
}

func InitDefaultRbac(uid int) {
	now := int(time.Now().Unix())
	roles := []rbac.Role{
		{Name: "superadmin", Description: "超级管理员", CreatedAt: now, UpdatedAt: now},
		{Name: "admin", Description: "管理员", CreatedAt: now, UpdatedAt: now},
		{Name: "editor", Description: "编辑", CreatedAt: now, UpdatedAt: now},
		{Name: "author", Description: "作者", CreatedAt: now, UpdatedAt: now},
		{Name: "contributor", Description: "贡献者", CreatedAt: now, UpdatedAt: now},
		{Name: "subscriber", Description: "订阅者", CreatedAt: now, UpdatedAt: now},
	}
	for _, role := range roles {
		db.GetDB().FirstOrCreate(&role, rbac.Role{Name: role.Name})
	}

	permissions := []rbac.Permission{
		{Name: "manage_users", Description: "管理用户", Category: "user", CreatedAt: now},
		{Name: "manage_content", Description: "管理内容", Category: "content", CreatedAt: now},
		{Name: "manage_comments", Description: "管理评论", Category: "comment", CreatedAt: now},
		{Name: "manage_meta", Description: "管理分类标签", Category: "meta", CreatedAt: now},
		{Name: "manage_system", Description: "管理系统", Category: "system", CreatedAt: now},
	}
	for _, perm := range permissions {
		db.GetDB().FirstOrCreate(&perm, rbac.Permission{Name: perm.Name})
	}

	var superadminRole rbac.Role
	db.GetDB().Where("name = ?", "superadmin").First(&superadminRole)
	userRole := rbac.UserRole{UserId: uid, RoleId: superadminRole.Id, GrantedAt: now, GrantedBy: 0}
	db.GetDB().FirstOrCreate(&userRole, rbac.UserRole{UserId: uid, RoleId: superadminRole.Id})

	var perms []rbac.Permission
	db.GetDB().Find(&perms)
	for _, perm := range perms {
		rolePerm := rbac.RolePermission{RoleId: superadminRole.Id, PermissionId: perm.Id}
		db.GetDB().FirstOrCreate(&rolePerm, rbac.RolePermission{RoleId: superadminRole.Id, PermissionId: perm.Id})
	}
}

func createDefaultUsermeta(uid int) {
	now := int(time.Now().Unix())
	defaults := []struct{ key, val string }{
		{"gender", "0"}, {"birthday", ""}, {"bio", "系统默认签名，送给每一位小可爱~"},
	}
	for _, d := range defaults {
		db.GetDB().Create(&models.Usermeta{UserId: uid, MetaKey: d.key, MetaValue: d.val, CreatedAt: now, UpdatedAt: now})
	}
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
