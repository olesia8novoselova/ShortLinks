package storage

import (
	"short-links/internal/models"
	"sync"
	"fmt"
)

type MemoryStorage struct {
	urls map[string]string
	mu   sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls: make(map[string]string),
	}
}

func (s *MemoryStorage) Save(url models.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls[url.Short] = url.Original
	return nil
}

func (s *MemoryStorage) Get(shortURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	originalURL, exists := s.urls[shortURL]
	if !exists {
		return "", fmt.Errorf("URL not found")
	}
	return originalURL, nil
}

func (s *MemoryStorage) Exists(shortURL string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.urls[shortURL]
	return exists
}