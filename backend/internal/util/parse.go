package util

import (
	"strconv"
)

// parseIntOrDefault attempts to parse a string to an int64.
// If parsing fails, it returns the provided defaultValue.
func ParseIntOrDefault(s string, defaultValue int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

// parseInt32OrDefault attempts to parse a string to an int32.
// If parsing fails, it returns the provided defaultValue.
func ParseInt32OrDefault(s string, defaultValue int32) int32 {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return int32(val)
}

// parseIntOrDefault attempts to parse a string to a boolean value
// If parsing fails, it returns the provided defaultValue.
func ParseBoolOrDefault(s string, defaultValue bool) bool {
	val, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue
	}
	return val
}
