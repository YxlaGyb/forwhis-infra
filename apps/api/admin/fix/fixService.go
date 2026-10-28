package fix

import (
	"fmt"

	"forestwhisperGo/apps/api/internal/models"
	"forestwhisperGo/apps/api/pkg/db"
)

type TableStatus struct {
	TableName string `json:"tableName"`
	Exists    bool   `json:"exists"`
	Status    string `json:"status"`
}

type tabler interface {
	TableName() string
}

func modelTableName(model any) string {
	if t, ok := model.(tabler); ok {
		return t.TableName()
	}
	return fmt.Sprintf("%T", model)
}

func registry() map[string]any {
	result := make(map[string]any)
	for _, model := range models.All() {
		result[modelTableName(model)] = model
	}
	return result
}

func CheckTables() []TableStatus {
	database := db.GetDB()
	if database == nil {
		return nil
	}

	tables := make([]TableStatus, 0, len(models.All()))
	for _, model := range models.All() {
		name := modelTableName(model)
		status := TableStatus{TableName: name}
		if database.Migrator().HasTable(model) {
			status.Exists = true
			status.Status = "正常"
		} else {
			status.Status = "缺失"
		}
		tables = append(tables, status)
	}
	return tables
}

func FixTables() (int, []string) {
	database := db.GetDB()
	if database == nil {
		return 0, []string{"数据库未连接"}
	}

	all := models.All()
	if err := database.AutoMigrate(all...); err != nil {
		return 0, []string{fmt.Sprintf("同步数据库结构失败: %v", err)}
	}
	return len(all), []string{"数据库结构已同步"}
}

func FixSingleTable(tableName string) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("数据库未连接")
	}

	model, ok := registry()[tableName]
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

	if _, ok := registry()[tableName]; !ok {
		return nil, fmt.Errorf("未知的表: %s", tableName)
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
	rows, err := database.Raw("SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = ?", tableName).Rows()
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
