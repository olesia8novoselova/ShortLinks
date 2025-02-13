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

	// инициализация PostgreSQL хранилища
	store, err := NewPostgresStorage(connStr)
	if err != nil {
		t.Fatalf("Не удалось инициализировать PostgreSQL хранилище: %v", err)
	}

	// тест Save и Get
	url := models.URL{Original: "https://example.com", Short: "abc123ABC_"}
	if err := store.Save(url); err != nil {
		t.Errorf("Не удалось сохранить URL: %v", err)
	}

	// тест получения сохраненного URL
	original, err := store.Get(url.Short)
	if err != nil {
		t.Errorf("Не удалось получить URL: %v", err)
	}
	if original != url.Original {
		t.Errorf("Ожидался оригинальный URL %s, получен %s", url.Original, original)
	}

	// тест несуществующего URL
	_, err = store.Get("invalid")
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующего URL")
	}

	// тест Exists
	if !store.Exists(url.Short) {
		t.Error("Ожидалось, что URL существует")
	}
	if store.Exists("invalid") {
		t.Error("Ожидалось, что URL не существует")
	}
}
