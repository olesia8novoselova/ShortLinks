package storage

import (
	"short-links/internal/models"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	store := NewMemoryStorage()

	// тест Save и Get
	url := models.URL{Original: "https://example.com", Short: "abc123ABC_"}
	if err := store.Save(url); err != nil {
		t.Errorf("Не удалось сохранить URL: %v", err)
	}

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
