package storage

import (
	"fmt"
	"short-links/internal/config"
	"short-links/internal/models"
	"testing"
)

func TestPostgresStorage(t *testing.T) {
	dbConfig := config.LoadDBConfig()
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.SSLMode,
	)

	// initialize PostgreSQL storage
	store, err := NewPostgresStorage(connStr)
	if err != nil {
		t.Fatalf("Failed to initialize PostgreSQL storage: %v", err)
	}

	// test Save and Get
	url := models.URL{Original: "https://example.com", Short: "abc123ABC_"}
	if err := store.Save(url); err != nil {
		t.Errorf("Failed to save URL: %v", err)
	}

	// test retrieving saved URL
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
