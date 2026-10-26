package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
	DBTimeZone string

	JWTSecret string
	AppPort   string

	UploadDir      string
	MaxUploadBytes int64
}

func LoadConfig() *Config {
	err := godotenv.Overload()
	if err != nil {
		log.Println("No .env file found, reading from system environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "amami_db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		DBTimeZone: getEnv("DB_TIMEZONE", "Asia/Jakarta"),

		JWTSecret: getEnv("JWT_SECRET", "super-secret-key-change-me"),
		AppPort:   getEnv("APP_PORT", "8080"),

		UploadDir:      getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadBytes: getEnvInt64("MAX_UPLOAD_MB", 10) << 20,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		log.Printf("invalid %s=%q, using default %d", key, value, defaultValue)
		return defaultValue
	}
	return parsed
}
