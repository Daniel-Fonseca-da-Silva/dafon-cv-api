package config

import (
	"os"
	"strconv"
)

// ParseIntEnv reads an integer from the environment, or returns defaultVal if missing/invalid.
func ParseIntEnv(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}

// ParseFloatEnv reads a float from the environment, or returns defaultVal if missing/invalid.
func ParseFloatEnv(key string, defaultVal float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultVal
	}
	return n
}
