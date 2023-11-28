package internal

import (
	"os"
	"strconv"
	"strings"
)

// GetAsString reads the environment variable or returns the default value
func GetAsString(key, dValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return dValue
}

// GetAsBool reads the environment variable into a boolean or returns the default value
func GetAsBool(name string, dValue bool) bool {
	str := GetAsString(name, "")
	if v, e := strconv.ParseBool(str); e == nil {
		return v
	}

	return dValue
}

// GetAsInt reads environment variable into an integer or returns the default value
func GetAsInt(name string, dValue int) int {
	str := GetAsString(name, "")
	if v, e := strconv.Atoi(str); e == nil {
		return v
	}

	return dValue
}

// GetAsSlice reads environment variable into a string slice or returns the default value
func GetAsSlice(name string, dValue []string, sep string) []string {
	str := GetAsString(name, "")
	if str == "" {
		return dValue
	}

	return strings.Split(str, sep)
}
