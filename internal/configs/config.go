package configs

import (
	"errors"
	"os"
)

type Config struct {
	Port  string
	DBUrl string
	Env   string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:  getEnv("APP_PORT", "3000"),
		DBUrl: getEnv("DB_URL", ""),
		Env:   getEnv("ENVIRONMENT", "development"),
	}

	if cfg.DBUrl == "" {
		return nil, errors.New("DB_URL environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); !ok {
		return fallback
	} else {
		return val
	}
}
