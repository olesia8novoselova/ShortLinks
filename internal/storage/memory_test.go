package storage

import (
	"short-links/internal/models"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	store := NewMemoryStorage()

	// test Save and Get
	url := models.URL{Original: "https://example.com", Short: "abc123ABC_"}
	if err := store.Save(url); err != nil {
		t.Errorf("Failed to save URL: %v", err)
	}

	original, err := store.Get(url.Short)
	if err != nil {
		t.Errorf("Failed to get URL: %v", err)
	}
	if original != url.Original {
		t.Errorf("Expected original URL %s, got %s", url.Original, original)
	}

	// test non-existent URL
	_, err = store.Get("invalid")
	if err == nil {
		t.Error("Expected error for non-existent URL")
	}

	// test Exists
	if !store.Exists(url.Short) {
		t.Error("Expected URL to exist")
	}
	if store.Exists("invalid") {
		t.Error("Expected URL to not exist")
	}
}
