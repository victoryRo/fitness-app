package main

import (
	"os"
	"strconv"
)

// GetAsString lee una variable de entorno o devuelve un valor por default
func GetAsString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetAsInt lee una variable de entorno como un numero o devuelve un valor por default
func GetAsInt(name string, defaultValue int) int {
	valueStr := GetAsString(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
