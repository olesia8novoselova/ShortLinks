package tests

import (
	"short-links/internal/handlers"
	"short-links/internal/storage"
	"short-links/internal/utils"
	"testing"
	"strings"
	"net/http"
	"net/http/httptest"
)

func TestGenerateShortUrl(t *testing.T) {
	original_url := "https://example.com"
	short_url := utils.GenerateShortUrl(original_url)

	if len(short_url) != 10 {
		t.Error("Expected short URL length to be 10, but got %d", len(short_url))
	}
}

func TestShortenUrl(t *testing.T) {
	urlStorage := storage.NewMemoryStorage()
	handler := handlers.NewHandler(urlStorage)

	requestBody := `{"original_url": "https://example.com"}`
	request := httptest.NewRequest("POST", "/shorten", strings.NewReader(requestBody))
	w := httptest.NewRecorder()

	handler.ShortenURL(w, request)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
	}
}