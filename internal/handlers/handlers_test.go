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
		t.Error("Ожидаемая длина короткого URL должна быть 10, но получено ", len(short_url))
	}
}

func TestShortenUrl(t *testing.T) {
	urlStorage := storage.NewMemoryStorage()
	handler := NewHandler(urlStorage)

	// тест валидного запроса
	requestBody := `{"original_url": "https://example.com"}`
	request := httptest.NewRequest("POST", "/shorten", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ShortenURL(w, request)

	if w.Code != http.StatusCreated {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusCreated, w.Code)
	}

	// тест невалидного запроса (отсутствует original_url)
	requestBody = `{"invalid": "data"}`
	request = httptest.NewRequest("POST", "/shorten", strings.NewReader(requestBody))
	w = httptest.NewRecorder()

	handler.ShortenURL(w, request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusBadRequest, w.Code)
	}

	// тест невалидного JSON
	requestBody = `invalid json`
	request = httptest.NewRequest("POST", "/shorten", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.ShortenURL(w, request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusBadRequest, w.Code)
	}

	// тест некорректного формата URL
	requestBody = `{"original_url": "invalid-url"}`
	request = httptest.NewRequest("POST", "/shorten", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.ShortenURL(w, request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusBadRequest, w.Code)
	}

	// тест пустого original_url
	requestBody = `{"original_url": ""}`
	request = httptest.NewRequest("POST", "/shorten", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.ShortenURL(w, request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetOriginalUrl(t *testing.T) {
	urlStorage := storage.NewMemoryStorage()
	handler := NewHandler(urlStorage)

	if err := urlStorage.Save(models.URL{Original: "https://example.com", Short: "abc123ABC_"}); err != nil {
		t.Errorf("Не удалось сохранить URL: %v", err)
	}

	// тест валидного запроса
	request := httptest.NewRequest("GET", "/original?short_url=abc123ABC_", nil)
	w := httptest.NewRecorder()

	handler.GetOriginalURL(w, request)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusOK, w.Code)
	}

	// тест отсутствия параметра short_url
	request = httptest.NewRequest("GET", "/original", nil)
	w = httptest.NewRecorder()

	handler.GetOriginalURL(w, request)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusBadRequest, w.Code)
	}

	// тест несуществующего short_url
	request = httptest.NewRequest("GET", "/original?short_url=nonexistent", nil)
	w = httptest.NewRecorder()

	handler.GetOriginalURL(w, request)

	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusNotFound, w.Code)
	}
}
