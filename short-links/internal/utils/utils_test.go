package utils

import (
	"testing"
)

func TestGenerateShortUrl(t *testing.T) {
	original := "https://example.com"
	short := GenerateShortUrl(original)

	if len(short) != 10 {
		t.Errorf("Expected short URL length 10, got %d", len(short))
	}
}

func TestGenerateRandomSymbol(t *testing.T) {
	symbol := GenerateRandomSymbol()
	if len(symbol) != 1 {
		t.Errorf("Expected symbol length 1, got %d", len(symbol))
	}
}