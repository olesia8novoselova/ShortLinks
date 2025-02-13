package storage

import (
	"fmt"
	"short-links/internal/models"
	"sync"
)

// map для хранения URL и mutex для потокобезопасности
type MemoryStorage struct {
	urls map[string]string
	mu   sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls: make(map[string]string),
	}
}

// Save сохраняет URL в памяти
func (s *MemoryStorage) Save(url models.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls[url.Short] = url.Original
	return nil
}

// Get извлекает оригинальный URL из хранилища по его сокращенной версии
func (s *MemoryStorage) Get(shortURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	originalURL, exists := s.urls[shortURL]
	if !exists {
		return "", fmt.Errorf("URL не найден")
	}
	return originalURL, nil
}

// Exists проверяет, существует ли сокращенный URL в хранилище
func (s *MemoryStorage) Exists(shortURL string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.urls[shortURL]
	return exists
}
