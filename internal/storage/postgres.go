package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"short-links/internal/models"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("Failed to create table: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

func createTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS urls (
		original_url TEXT PRIMARY KEY,
		short_url TEXT UNIQUE
	);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("Failed to create table: %w", err)
	}
	return nil
}

func (s *PostgresStorage) Save(url models.URL) error {
	query := `INSERT INTO urls (original_url, short_url) 
	VALUES ($1, $2)
	ON CONFLICT (original_url) DO NOTHING;`
	_, err := s.db.Exec(query, url.Original, url.Short)
	if err != nil {
		return fmt.Errorf("Failed to save URL: %w", err)
	}
	return nil
}

func (s *PostgresStorage) Get(short_url string) (string, error){
	var original_url string
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_url = $1);`
	err := s.db.QueryRow(query, short_url).Scan(&original_url)
	if errors.Is(err, sql.ErrNoRows) {
		return "Shorten URL not found", nil
	}
	return original_url, err
}

func (s *PostgresStorage) Exists(short_url string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_url = $1);`
	_ = s.db.QueryRow(query, short_url).Scan(&exists)
	return exists
}