package storage

import (
	"short-links/internal/models"
	"sync"
	"fmt"
)

// map of URLs and a mutex for thread safety
type MemoryStorage struct {
	urls map[string]string
	mu   sync.Mutex
}


func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls: make(map[string]string),
	}
}

// Save saves a URL in the memory storage
func (s *MemoryStorage) Save(url models.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls[url.Short] = url.Original
	return nil
}

// Get retrieves the original URL from the storage by its shortened version
func (s *MemoryStorage) Get(shortURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	originalURL, exists := s.urls[shortURL]
	if !exists {
		return "", fmt.Errorf("URL not found")
	}
	return originalURL, nil
}

// Exists checks if a shortened URL exists in the storage
func (s *MemoryStorage) Exists(shortURL string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.urls[shortURL]
	return exists
}