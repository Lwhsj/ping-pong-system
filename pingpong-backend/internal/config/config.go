package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	AppEnv      string
	DBHost      string
	DBPort      string
	DBName      string
	DBUser      string
	DBPassword  string
	DBParseTime string
	DBLoc       string
	UploadDir   string
	MaxUploadMB int64
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "development"),
		DBHost:      getEnv("DB_HOST", "127.0.0.1"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBName:      getEnv("DB_NAME", "pingpong"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", "1234"),
		DBParseTime: getEnv("DB_PARSE_TIME", "true"),
		DBLoc:       getEnv("DB_LOC", "Local"),
		UploadDir:   getEnv("UPLOAD_DIR", "uploads"),
		MaxUploadMB: getEnvInt64("MAX_UPLOAD_MB", 50),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=%s&loc=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBParseTime,
		c.DBLoc,
	)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
