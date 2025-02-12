package handlers

import (
	"encoding/json"
	"net/http"
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

// POST request
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var url models.URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// ensure original_url is not empty
	if url.Original == "" {
		http.Error(w, "Missing or empty original_url", http.StatusBadRequest)
		return
	}

	url.Short = utils.GenerateShortUrl(url.Original)

	// ensure unique short URL
	for h.storage.Exists(url.Short) {
		url.Short += utils.GenerateRandomSymbol()
	}

	if err := h.storage.Save(url); err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(url)
}


// GET request
func (h *Handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Query().Get("short_url")
	if shortURL == "" {
		http.Error(w, "Missing short_url parameter", http.StatusBadRequest)
		return
	}

	// ensure short URL exists
	originalURL, err := h.storage.Get(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"original_url": originalURL})
}