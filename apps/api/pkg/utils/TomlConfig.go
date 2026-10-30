package utils

import (
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

var (
	configPath = "config.toml"
	config     map[string]map[string]interface{}
	configOnce sync.Once
)

func LoadConfig() {
	data, err := os.ReadFile(configPath)
	if err != nil {
		config = make(map[string]map[string]interface{})
		return
	}
	toml.Unmarshal(data, &config)
}

func GetConfig(section, key string) interface{} {
	configOnce.Do(LoadConfig)
	if s, ok := config[section]; ok {
		return s[key]
	}
	return nil
}

func GetString(section, key string) string {
	if v := GetConfig(section, key); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func GetInt(section, key string) int {
	if v := GetConfig(section, key); v != nil {
		switch n := v.(type) {
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
