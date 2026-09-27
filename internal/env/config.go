package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBHost     string
	DBName     string
	DBPort     string
	DBPassword string
	DBSecure   string
	JWTSecret  string
	ServerPort string
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func findEnvFile() string {
	dir, _ := os.Getwd()
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func Load() (*Config, error) {
	if path := findEnvFile(); path != "" {
		_ = godotenv.Load(path)
	}

	cfg := &Config{
		getEnv("DB_USER", " "),
		getEnv("DB_HOST", " "),
		getEnv("DB_NAME", " "),
		getEnv("DB_PORT", " "),
		getEnv("DB_PASS", " "),
		getEnv("DB_SECURE", " "),
		getEnv("TEST_JWT_SECRET", " "),
		getEnv("SERVER_PORT", " "),
	}

	if cfg.DBUser == "" {
		return nil, fmt.Errorf("DB_USER missing")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD  Missing")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("DB_NAME  Missing")
	}

	return cfg, nil
}
