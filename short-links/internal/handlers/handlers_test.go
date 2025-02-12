package handlers

import (
	"net/http"
	"net/http/httptest"
	"short-links/internal/models"
	"short-links/internal/storage"
	"short-links/internal/utils"
	"strings"
	"testing"
)

func TestGenerateShortUrl(t *testing.T) {
	original_url := "https://example.com"
	short_url := utils.GenerateShortUrl(original_url)

	if len(short_url) != 10 {
		t.Error("Expected short URL length to be 10, but got ", len(short_url))
	}
}

func TestShortenUrl(t *testing.T) {
	urlStorage := storage.NewMemoryStorage()
	handler := NewHandler(urlStorage)

	// test valid request
	requestBody := `{"original_url": "https://example.com"}`
	request := httptest.NewRequest("POST", "/shorten", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ShortenURL(w, request)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

	// test invalid request (missing original_url)
	requestBody = `{"invalid": "data"}`
	request = httptest.NewRequest("POST", "/shorten", strings.NewReader(requestBody))
	w = httptest.NewRecorder()

	handler.ShortenURL(w, request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}

	// test invalid JSON
	requestBody = `invalid json`
	request = httptest.NewRequest("POST", "/shorten", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.ShortenURL(w, request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetOriginalUrl(t *testing.T) {
	urlStorage := storage.NewMemoryStorage()
	handler := NewHandler(urlStorage)

	if err := urlStorage.Save(models.URL{Original: "https://example.com", Short: "abc123ABC_"}); err != nil {
		t.Errorf("Failed to save URL: %v", err)
	}

	// test valid request
	request := httptest.NewRequest("GET", "/original?short_url=abc123ABC_", nil)
	w := httptest.NewRecorder()

	handler.GetOriginalURL(w, request)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	// test missing short_url parameter
	request = httptest.NewRequest("GET", "/original", nil)
	w = httptest.NewRecorder()

	handler.GetOriginalURL(w, request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}

	// test non-existing short_url
	request = httptest.NewRequest("GET", "/original?short_url=nonexistent", nil)
	w = httptest.NewRecorder()

	handler.GetOriginalURL(w, request)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, w.Code)
	}
}
