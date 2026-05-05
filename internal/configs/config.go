package configs

import "os"

type Config struct {
	Port  string
	DbUrl string
}

func Load() *Config {
	return &Config{
		Port:  getEnv("APP_PORT", "3000"),
		DbUrl: getEnv("DB_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); !ok {
		return fallback
	} else {
		return val
	}
}
