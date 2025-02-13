package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"short-links/internal/config"
	"short-links/internal/handlers"
	"short-links/internal/storage"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	// тип хранилища указывается как параметр при запуске сервиса
	storageType := flag.String("storage", "memory", "Тип хранилища (memory|postgres)")
	flag.Parse()

	var store storage.Storage

	switch *storageType {
	case "postgres":
		dbConfig := config.LoadDBConfig()
		connStr := formatConnectionString(dbConfig)
		store, err = storage.NewPostgresStorage(connStr)
	case "memory":
		store = storage.NewMemoryStorage()
	default:
		log.Fatalf("Неверный тип хранилища: %s", *storageType)
	}

	if err != nil {
		log.Fatalf("Не удалось инициализировать хранилище: %v", err)
	}

	handler := handlers.NewHandler(store)

	http.HandleFunc("/shorten", handler.ShortenURL)
	http.HandleFunc("/original", handler.GetOriginalURL)

	port := getPort()
	log.Printf("Сервер запущен на порту %s с использованием хранилища %s", port, *storageType)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func formatConnectionString(cfg *config.DBConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode)
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}
