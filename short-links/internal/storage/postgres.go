package storage

import (
	"database/sql"
	"fmt"
	"short-links/internal/models"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	// open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// check if the database is reachable
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// create the URLs table if it doesn't exist
	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// createTable ensures the 'urls' table exists in the database
func createTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS urls (
		original_url TEXT PRIMARY KEY,
		short_url TEXT UNIQUE
	);`
	_, err := db.Exec(query)
	return err
}

// Save stores a URL in the PostgreSQL database
func (s *PostgresStorage) Save(url models.URL) error {
	query := `INSERT INTO urls (original_url, short_url) 
	VALUES ($1, $2)
	ON CONFLICT (original_url) DO NOTHING;`
	_, err := s.db.Exec(query, url.Original, url.Short)
	return err
}

// Get retrieves the original URL from the database using the short URL
func (s *PostgresStorage) Get(shortURL string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_url = $1`
	err := s.db.QueryRow(query, shortURL).Scan(&originalURL)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("short URL not found")
	}
	return originalURL, err
}

// Exists checks if a short URL exists in the database
func (s *PostgresStorage) Exists(shortURL string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_url = $1)`
	err := s.db.QueryRow(query, shortURL).Scan(&exists)
	return err == nil && exists
}