package service

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/pkg/db"
)

var (
	optionsCache     = make(map[string]string)
	optionsCacheLock sync.RWMutex
	optionsLoaded    bool
)

func GetOption(key string, defaultValue string) string {
	optionsCacheLock.RLock()
	if val, ok := optionsCache[key]; ok {
		optionsCacheLock.RUnlock()
		return val
	}
	optionsCacheLock.RUnlock()

	var option models.Option
	if db.GetDB().Where("option_key = ?", key).First(&option).Error != nil {
		return defaultValue
	}

	optionsCacheLock.Lock()
	optionsCache[option.OptionKey] = option.OptionValue
	optionsCacheLock.Unlock()

	return option.OptionValue
}

func GetOptionInt(key string, defaultValue int) int {
	val := GetOption(key, "")
	if val == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func GetOptionBool(key string, defaultValue bool) bool {
	val := GetOption(key, "")
	if val == "" {
		return defaultValue
	}
	return val == "1" || val == "true" || val == "yes"
}

func GetOptionJSON(key string, target interface{}) error {
	val := GetOption(key, "")
	if val == "" {
		return nil
	}
	return json.Unmarshal([]byte(val), target)
}

func SetOption(key, value string) error {
	now := int(time.Now().Unix())
	var option models.Option
	result := db.GetDB().Where("option_key = ?", key).First(&option)
	if result.Error != nil {
		option = models.Option{
			OptionKey:   key,
			OptionValue: value,
			Autoload:    models.AutoloadNo,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := db.GetDB().Create(&option).Error; err != nil {
			return err
		}
	} else {
		option.OptionValue = value
		option.UpdatedAt = now
		if err := db.GetDB().Save(&option).Error; err != nil {
			return err
		}
	}

	optionsCacheLock.Lock()
	optionsCache[key] = value
	optionsCacheLock.Unlock()

	return nil
}

func SetOptionInt(key string, value int) error {
	return SetOption(key, strconv.Itoa(value))
}

func SetOptionBool(key string, value bool) error {
	if value {
		return SetOption(key, "1")
	}
	return SetOption(key, "0")
}

func SetOptionJSON(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return SetOption(key, string(data))
}

func DeleteOption(key string) bool {
	result := db.GetDB().Where("option_key = ?", key).Delete(&models.Option{})
	if result.Error != nil {
		return false
	}

	optionsCacheLock.Lock()
	delete(optionsCache, key)
	optionsCacheLock.Unlock()

	return result.RowsAffected > 0
}

func LoadAutoloadOptions() error {
	var options []models.Option
	if err := db.GetDB().Where("autoload = ?", models.AutoloadYes).Find(&options).Error; err != nil {
		return err
	}

	optionsCacheLock.Lock()
	for _, opt := range options {
		optionsCache[opt.OptionKey] = opt.OptionValue
	}
	optionsLoaded = true
	optionsCacheLock.Unlock()

	return nil
}

func GetAllOptions() map[string]string {
	var options []models.Option
	db.GetDB().Find(&options)

	result := make(map[string]string)
	for _, opt := range options {
		result[opt.OptionKey] = opt.OptionValue
	}
	return result
}

func GetOptionsByPrefix(prefix string) map[string]string {
	var options []models.Option
	db.GetDB().Where("option_key LIKE ?", prefix+"%").Find(&options)

	result := make(map[string]string)
	for _, opt := range options {
		result[opt.OptionKey] = opt.OptionValue
	}
	return result
}

func SetOptions(options map[string]string) error {
	tx := db.GetDB().Begin()
	now := int(time.Now().Unix())

	for key, value := range options {
		var option models.Option
		result := tx.Where("option_key = ?", key).First(&option)
		if result.Error != nil {
			option = models.Option{
				OptionKey:   key,
				OptionValue: value,
				Autoload:    models.AutoloadNo,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			if err := tx.Create(&option).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			option.OptionValue = value
			option.UpdatedAt = now
			if err := tx.Save(&option).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	optionsCacheLock.Lock()
	for key, value := range options {
		optionsCache[key] = value
	}
	optionsCacheLock.Unlock()

	return nil
}

func ClearOptionsCache() {
	optionsCacheLock.Lock()
	optionsCache = make(map[string]string)
	optionsLoaded = false
	optionsCacheLock.Unlock()
}
