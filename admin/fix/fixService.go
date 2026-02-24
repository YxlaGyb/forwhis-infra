package fix

import (
	"fmt"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/rbac"
	"forestwhisperGo/pkg/db"
)

type TableStatus struct {
	TableName string `json:"tableName"`
	Exists    bool   `json:"exists"`
	Status    string `json:"status"`
}

func CheckTables() []TableStatus {
	database := db.GetDB()
	if database == nil {
		return nil
	}

	tables := []TableStatus{
		{TableName: "users"},
		{TableName: "usermetas"},
		{TableName: "userinfos"},
		{TableName: "contents"},
		{TableName: "comments"},
		{TableName: "metas"},
		{TableName: "relationships"},
		{TableName: "reviews"},
		{TableName: "content_fields"},
		{TableName: "ai_conversations"},
		{TableName: "likes"},
		{TableName: "roles"},
		{TableName: "permissions"},
		{TableName: "user_roles"},
		{TableName: "role_permissions"},
	}

	for i := range tables {
		if database.Migrator().HasTable(tables[i].TableName) {
			tables[i].Exists = true
			tables[i].Status = "正常"
		} else {
			tables[i].Exists = false
			tables[i].Status = "缺失"
		}
	}

	return tables
}

func FixTables() (int, []string) {
	database := db.GetDB()
	if database == nil {
		return 0, []string{"数据库未连接"}
	}

	results := make([]string, 0)
	fixedCount := 0

	modelsToMigrate := []struct {
		name  string
		model interface{}
	}{
		{"users", &models.User{}},
		{"usermetas", &models.Usermeta{}},
		{"userinfos", &models.Userinfo{}},
		{"contents", &models.Content{}},
		{"comments", &models.Comment{}},
		{"metas", &models.Meta{}},
		{"relationships", &models.Relationship{}},
		{"reviews", &models.Review{}},
		{"content_fields", &models.ContentField{}},
		{"ai_conversations", &models.AIConversation{}},
		{"likes", &models.Like{}},
		{"roles", &rbac.Role{}},
		{"permissions", &rbac.Permission{}},
		{"user_roles", &rbac.UserRole{}},
		{"role_permissions", &rbac.RolePermission{}},
	}

	for _, m := range modelsToMigrate {
		if !database.Migrator().HasTable(m.name) {
			if err := database.AutoMigrate(m.model); err != nil {
				results = append(results, fmt.Sprintf("修复 %s 失败: %v", m.name, err))
			} else {
				results = append(results, fmt.Sprintf("修复 %s 成功", m.name))
				fixedCount++
			}
		}
	}

	if fixedCount == 0 {
		results = append(results, "所有表都正常, 无需修复")
	}

	return fixedCount, results
}

func FixSingleTable(tableName string) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("数据库未连接")
	}

	modelsMap := map[string]interface{}{
		"users":            &models.User{},
		"usermetas":        &models.Usermeta{},
		"userinfos":        &models.Userinfo{},
		"contents":         &models.Content{},
		"comments":         &models.Comment{},
		"metas":            &models.Meta{},
		"relationships":    &models.Relationship{},
		"reviews":          &models.Review{},
		"content_fields":   &models.ContentField{},
		"ai_conversations": &models.AIConversation{},
		"likes":            &models.Like{},
		"roles":            &rbac.Role{},
		"permissions":      &rbac.Permission{},
		"user_roles":       &rbac.UserRole{},
		"role_permissions": &rbac.RolePermission{},
	}

	model, ok := modelsMap[tableName]
	if !ok {
		return fmt.Errorf("未知的表: %s", tableName)
	}

	return database.AutoMigrate(model)
}

func GetTableInfo(tableName string) (map[string]interface{}, error) {
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	if !database.Migrator().HasTable(tableName) {
		return map[string]interface{}{
			"exists": false,
			"name":   tableName,
		}, nil
	}

	var count int64
	database.Table(tableName).Count(&count)

	columns := make([]map[string]interface{}, 0)
	rows, err := database.Raw(fmt.Sprintf("SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = '%s'", tableName)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName, dataType, isNullable string
		rows.Scan(&columnName, &dataType, &isNullable)
		columns = append(columns, map[string]interface{}{
			"name":     columnName,
			"type":     dataType,
			"nullable": isNullable,
		})
	}

	return map[string]interface{}{
		"exists":  true,
		"name":    tableName,
		"count":   count,
		"columns": columns,
	}, nil
}

func CheckFixToolAllowed() bool {
	return db.GetDB() != nil
}

func CreateMissingIndexes() (int, []string) {
	database := db.GetDB()
	if database == nil {
		return 0, []string{"数据库未连接"}
	}

	results := make([]string, 0)
	createdCount := 0

	indexes := []struct {
		table string
		name  string
		sql   string
	}{
		{"likes", "idx_likes_user_id", "CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id)"},
		{"likes", "idx_likes_target_type", "CREATE INDEX IF NOT EXISTS idx_likes_target_type ON likes(target_type)"},
		{"likes", "idx_likes_target_id", "CREATE INDEX IF NOT EXISTS idx_likes_target_id ON likes(target_id)"},
		{"likes", "idx_likes_unique", "CREATE UNIQUE INDEX IF NOT EXISTS idx_likes_unique ON likes(user_id, target_type, target_id)"},
	}

	for _, idx := range indexes {
		if !database.Migrator().HasTable(idx.table) {
			continue
		}

		if err := database.Exec(idx.sql).Error; err != nil {
			results = append(results, fmt.Sprintf("创建索引 %s 失败: %v", idx.name, err))
		} else {
			results = append(results, fmt.Sprintf("创建索引 %s 成功", idx.name))
			createdCount++
		}
	}

	if createdCount == 0 {
		results = append(results, "所有索引都正常, 无需创建")
	}

	return createdCount, results
}
