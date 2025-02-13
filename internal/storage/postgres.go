package storage

import (
	"database/sql"
	"fmt"
	"short-links/internal/models"

	_ "github.com/lib/pq" // драйвер PostgreSQL
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	// открываем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	// проверяем, доступна ли база данных
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось проверить соединение с базой данных: %w", err)
	}

	// создаем таблицу URLs, если она не существует
	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("не удалось создать таблицу: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// createTable гарантирует, что таблица 'urls' существует в базе данных
func createTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS urls (
		original_url TEXT PRIMARY KEY,
		short_url TEXT UNIQUE
	);`
	_, err := db.Exec(query)
	return err
}

// Save сохраняет URL в базе данных PostgreSQL
func (s *PostgresStorage) Save(url models.URL) error {
	query := `INSERT INTO urls (original_url, short_url) 
	VALUES ($1, $2)
	ON CONFLICT (original_url) DO NOTHING;`
	_, err := s.db.Exec(query, url.Original, url.Short)
	return err
}

// Get извлекает оригинальный URL из базы данных по короткому URL
func (s *PostgresStorage) Get(shortURL string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_url = $1`
	err := s.db.QueryRow(query, shortURL).Scan(&originalURL)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("короткий URL не найден")
	}
	return originalURL, err
}

// Exists проверяет, существует ли короткий URL в базе данных
func (s *PostgresStorage) Exists(shortURL string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_url = $1)`
	err := s.db.QueryRow(query, shortURL).Scan(&exists)
	return err == nil && exists
}
