package main

import (
	"os"
	"strconv"
)

// GetAsString lee una variable de entorno o devuelve un valor predeterminado
func GetAsString(key, defaultValue string) string {
	// LookupEnv recupera el valor env
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

// GetAsInt lee una variable de entorno en un número entero o devuelve un valor predeterminado
func GetAsInt(name string, defaultValue int) int {
	valueStr := GetAsString(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}
