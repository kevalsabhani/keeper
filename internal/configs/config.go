package configs

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all runtime configuration for the application.
// Values are populated from environment variables with sensible defaults.
type Config struct {
	Port         string
	DBUrl        string
	Env          string
	ReadTimeOut  int
	WriteTimeOut int
	IdleTimeOut  int
}

// Load reads configuration from environment variables and returns a validated Config.
// It returns an error if any required variable is missing or any value is malformed.
func Load() (*Config, error) {
	readTimeout, err := getEnvInt("READ_TIMEOUT", "15")
	if err != nil {
		return nil, err
	}

	writeTimeout, err := getEnvInt("WRITE_TIMEOUT", "30")
	if err != nil {
		return nil, err
	}

	idleTimeout, err := getEnvInt("IDLE_TIMEOUT", "30")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Port:         getEnv("APP_PORT", "3000"),
		DBUrl:        getEnv("DB_URL", ""),
		Env:          getEnv("ENVIRONMENT", "development"),
		ReadTimeOut:  readTimeout,
		WriteTimeOut: writeTimeout,
		IdleTimeOut:  idleTimeout,
	}

	if cfg.DBUrl == "" {
		return nil, fmt.Errorf("DB_URL environment variable is required")
	}

	return cfg, nil
}

// getEnvInt reads an environment variable as a base-10 integer.
// Returns an error if the resolved value cannot be parsed as an integer.
func getEnvInt(key, fallback string) (int, error) {
	valueStr := getEnv(key, fallback)
	seconds, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %q is not a valid integer", key, valueStr)
	}
	return seconds, nil
}

// getEnv returns the value of the environment variable named by key,
// or fallback if the variable is not set.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
