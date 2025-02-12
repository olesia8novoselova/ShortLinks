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
)

func main() {
	
	// the storage type is specified as a parameter when starting the service
	storageType := flag.String("storage", "memory", "Storage type (memory|postgres)")
	flag.Parse()

	var store storage.Storage
	var err error

	switch *storageType {
	case "postgres":
		dbConfig := config.LoadDBConfig()
		connStr := formatConnectionString(dbConfig)
		store, err = storage.NewPostgresStorage(connStr)
	case "memory":
		store = storage.NewMemoryStorage()
	default:
		log.Fatalf("Invalid storage type: %s", *storageType)
	}

	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	
	handler := handlers.NewHandler(store)

	http.HandleFunc("/shorten", handler.ShortenURL)
	http.HandleFunc("/original", handler.GetOriginalURL)

	port := getPort()
	log.Printf("Server started on port %s using %s storage", port, *storageType)
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