package env

import (
	"os"
	"strconv"
)

func GetString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func GetInt(key string, defaultVal int) int {
	valueStr := GetString(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func GetBool(key string, defaultVal bool) bool {
	valStr := GetString(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
