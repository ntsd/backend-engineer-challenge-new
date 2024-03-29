package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config is environment config
type Config struct {
	AppPort     string
	PostgresURL string
	ScanWorkers int
}

// NewConfig returns environment config
func NewConfig() Config {
	return Config{
		AppPort:     defaultENV("APP_PORT", "8080"),
		PostgresURL: requiredENV("POSTGRES_URL"),
		ScanWorkers: strToInt(defaultENV("SCAN_WORKERS", "2")),
	}
}

// requiredENV returns environment variable value, panic if not found
func requiredENV(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal(fmt.Errorf("`%s` env is required", key))
	}
	return value
}

// defaultENV will return environment variable or default value if it's empty
func defaultENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}

// strToInt convert string to int, panic if not found
func strToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(fmt.Errorf("`%s` env value can not convert to int", str))
	}
	return i
}
