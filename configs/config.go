package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// LoadConfig загружает конфигурацию из .env файла или переменных окружения
func LoadConfig() (Config, error) {
	_ = godotenv.Load()
	requiredVars := []string{"DB_HOST", "DB_PORT", "DB_USERNAME", "DB_PASSWORD", "DB_NAME"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Printf("Warning: Environment variable %s is not set\n", v)
		}
	}

	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUsername: getEnv("DB_USERNAME", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "postgres"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
