package configutil

func GetString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func GetInt(m map[string]interface{}, key string) int {
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

func GetBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func GetSectionString(m map[string]map[string]interface{}, section, key string) string {
	if s, ok := m[section]; ok {
		return GetString(s, key)
	}
	return ""
}

func GetSectionInt(m map[string]map[string]interface{}, section, key string) int {
	if s, ok := m[section]; ok {
		return GetInt(s, key)
	}
	return 0
}

func GetSectionBool(m map[string]map[string]interface{}, section, key string) bool {
	if s, ok := m[section]; ok {
		return GetBool(s, key)
	}
	return false
}
