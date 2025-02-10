package storage

import (
	"short-links/models"
	"sync"
)

type MemoryStorage struct {
	urls map[string]string
	mu sync.Mutex
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

func (s *MemoryStorage) Get(short_url string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	original_url, ok := s.urls[short_url]
	if !ok {
		return "URL not found", nil
	}
	return original_url, nil
}