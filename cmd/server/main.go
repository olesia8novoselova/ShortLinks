package main

import (
	"short-links/internal/handlers"
	"short-links/internal/storage"
	"net/http"
	"log"
)

func main() {
	urlStorage := storage.NewMemoryStorage()
	handler := handlers.NewHandler(urlStorage)

	http.HandleFunc("/shorten", handler.ShortenURL)
	http.HandleFunc("/original", handler.GetOriginalURL)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}