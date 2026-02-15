package db

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

var (
	cfg     map[string]map[string]interface{}
	cfgOnce sync.Once
)

func getStr(section, key string) string {
	cfgOnce.Do(func() {
		data, _ := os.ReadFile("config.toml")
		toml.Unmarshal(data, &cfg)
	})
	if s, ok := cfg[section]; ok {
		if v, ok := s[key].(string); ok {
			return v
		}
	}
	return ""
}

func getInt(section, key string) int {
	cfgOnce.Do(func() {
		data, _ := os.ReadFile("config.toml")
		toml.Unmarshal(data, &cfg)
	})
	if s, ok := cfg[section]; ok {
		switch n := s[key].(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		}
	}
	return 0
}

func initDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		getStr("db", "PGSQLHOST"),
		getStr("db", "PGSQLUSER"),
		getStr("db", "PGSQLPWD"),
		getStr("db", "PGSQL_DB"),
		getInt("db", "PGSQLPORT"),
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("数据库连接失败: %v", err)
		DB = nil
		return
	}
	log.Println("数据库连接成功")
}

func GetDB() *gorm.DB {
	once.Do(initDB)
	return DB
}

func CloseDB() {
	if DB != nil {
		sqlDB, _ := DB.DB()
		sqlDB.Close()
	}
}
