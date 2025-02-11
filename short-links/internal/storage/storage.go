package storage

import "short-links/internal/models"

type Storage interface {
	Save(url models.URL) error
	Get(shortURL string) (string, error)
	Exists(shortURL string) bool
}