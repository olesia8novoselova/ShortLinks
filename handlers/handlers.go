package handlers

import (
	"encoding/json"
	"short-links/models"
	"short-links/storage"
	"short-links/utils"
	"net/http"
)

// handling HTTP requests for URL shortening
type Handler struct {
	urlStorage storage.Storage
}

func NewHandler(urlStorage storage.Storage) *Handler {
	return &Handler{urlStorage: urlStorage}
}

// POST requests
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var url models.URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	url.Short = utils.GenerateShortUrl(url.Original)

	// ensure the short URL is unique
	for {
		_, err := h.urlStorage.Get(url.Short)
		if err != nil {
			break
		}
		url.Short += utils.GenerateRandomSymbol()
	}

	// save the URL in the storage
	err = h.urlStorage.Save(url)
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	// return the short URL in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(url)
}

// GET requests
func (h *Handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	short_url := r.URL.Query().Get("short_url")
	if short_url == "" {
		http.Error(w, "Missing short URL parameter", http.StatusBadRequest)
		return
	}

	// get the original URL from the storage
	original_url, err := h.urlStorage.Get(short_url)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// return the original URL in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"original_url": original_url})
}