package db

import (
	"fmt"
	"log"
	"os"
	"sync"

	"forestwhisperGo/magic/utils/configutil"

	"github.com/pelletier/go-toml/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	mu    sync.RWMutex
	tried bool
)

var (
	cfg     map[string]map[string]interface{}
	cfgOnce sync.Once
)

func loadConfig() {
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return
	}
	toml.Unmarshal(data, &cfg)
}

func initDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		configutil.GetSectionString(cfg, "db", "PGSQLHOST"),
		configutil.GetSectionString(cfg, "db", "PGSQLUSER"),
		configutil.GetSectionString(cfg, "db", "PGSQLPWD"),
		configutil.GetSectionString(cfg, "db", "PGSQL_DB"),
		configutil.GetSectionInt(cfg, "db", "PGSQLPORT"),
	)
	var err error
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("数据库连接失败: %v", err)
		return
	}
	sqlDB, _ := conn.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	log.Println("数据库连接成功")
	mu.Lock()
	DB = conn
	mu.Unlock()
}

func pingAndReconnect() bool {
	if DB == nil {
		return false
	}
	sqlDB, err := DB.DB()
	if err != nil {
		initDB()
		return DB != nil
	}
	if err := sqlDB.Ping(); err != nil {
		log.Printf("数据库连接失效，尝试重连: %v", err)
		initDB()
		return DB != nil
	}
	return true
}

func GetDB() *gorm.DB {
	cfgOnce.Do(loadConfig)

	if DB != nil && !pingAndReconnect() {
		return nil
	}

	mu.RLock()
	if DB != nil {
		mu.RUnlock()
		return DB
	}
	mu.RUnlock()

	if !tried {
		tried = true
		initDB()
	}
	mu.RLock()
	defer mu.RUnlock()
	return DB
}

func CloseDB() {
	mu.Lock()
	defer mu.Unlock()
	if DB != nil {
		sqlDB, _ := DB.DB()
		sqlDB.Close()
		DB = nil
	}
}
