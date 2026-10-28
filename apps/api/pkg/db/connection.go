package db

import (
	"fmt"
	"log"
	"os"
	"sync"

	"forestwhisperGo/apps/api/pkg/utils/configutil"

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
		log.Printf("[FATAL] 数据库连接失败: %v", err)
		return
	}
	sqlDB, err := conn.DB()
	if err != nil {
		log.Printf("[FATAL] 获取SQL连接失败: %v", err)
		return
	}
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
		log.Printf("数据库连接失效, 尝试重连: %v", err)
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

func CloseDB() error {
	mu.Lock()
	defer mu.Unlock()
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
		DB = nil
	}
	return nil
}

func ReloadDB() error {
	mu.Lock()
	if DB != nil {
		sqlDB, _ := DB.DB()
		sqlDB.Close()
		DB = nil
	}
	tried = false
	mu.Unlock()

	cfgOnce = sync.Once{}
	cfg = nil

	cfgOnce.Do(loadConfig)
	if cfg == nil {
		return fmt.Errorf("配置加载失败")
	}

	initDB()
	if DB == nil {
		return fmt.Errorf("数据库连接失败")
	}
	return nil
}
