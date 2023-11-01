package main

import (
	"os"
	"strconv"
)

// GetAsString reads a environment variable and returns it as a string.
// If the environment variable is not set, it returns a default value.
func GetAsString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

// GetAsInt reads a environment variable and returns it as an integer.
func GetAsInt(name string, defaultValue int) int {
	value := GetAsString(name, "")
	if v, e := strconv.Atoi(value); nil != e {
		return v
	}

	return defaultValue
}
