package storage

import (
	"short-links/internal/models"
)

type Storage interface {
	Save(url models.URL) error
	Get(short_url string) (string, error)
}