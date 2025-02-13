package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"short-links/internal/models"
	"short-links/internal/storage"
	"short-links/internal/utils"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

// Regex для валидации URL
var urlRegex = regexp.MustCompile(`^(http|https)://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/[a-zA-Z0-9-._~:?#@!$&'()*+,;=]*)*$`)

// POST запрос
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var url models.URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	// проверка, что original_url не пустой
	if url.Original == "" {
		http.Error(w, "Отсутствует или пустой original_url", http.StatusBadRequest)
		return
	}

	// проверка на валидность URL
	if !urlRegex.MatchString(url.Original) {
		http.Error(w, "Некорректный формат URL", http.StatusBadRequest)
		return
	}

	url.Short = utils.GenerateShortUrl(url.Original)

	if err := h.storage.Save(url); err != nil {
		http.Error(w, "Не удалось сохранить URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(url); err != nil {
		http.Error(w, "Не удалось закодировать ответ", http.StatusInternalServerError)
		return
	}
}

// GET запрос
func (h *Handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Query().Get("short_url")
	if shortURL == "" {
		http.Error(w, "Отсутствует параметр short_url", http.StatusBadRequest)
		return
	}

	// проверка, что короткий URL существует
	originalURL, err := h.storage.Get(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"original_url": originalURL}); err != nil {
		http.Error(w, "Не удалось закодировать ответ", http.StatusInternalServerError)
		return
	}
}
