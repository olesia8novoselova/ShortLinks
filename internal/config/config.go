package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadDBConfig() *DBConfig {
	// несколько путей для загрузки файла .env
	envPaths := []string{
		filepath.Join(".env"),             // для запуска сервера
		filepath.Join("..", "..", ".env"), // для запуска тестов
	}

	var err error
	for _, path := range envPaths {
		err = godotenv.Load(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		fmt.Println("Ошибка загрузки файла .env", err)
	}

	return &DBConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "5432"),
		User:     GetEnv("DB_USER", "postgres"),
		Password: GetEnv("DB_PASSWORD", ""),
		DBName:   GetEnv("DB_NAME", "urlshortener"),
		SSLMode:  GetEnv("SSL_MODE", "disable"),
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
