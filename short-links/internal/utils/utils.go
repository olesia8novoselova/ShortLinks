package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"strings"
)

// generating a short URL based on the original URL using SHA1 hash and Base64 encoding
func GenerateShortUrl(original_url string) string {
	// calculate SHA1 hash of the original URL
	hash := sha1.Sum([]byte(original_url))
	// encode the hash using Base64 URL encoding
	short_url := base64.URLEncoding.EncodeToString(hash[:])
	// remove padding characters
	short_url = strings.TrimRight(short_url, "=")

	if len(short_url) > 10 {
		short_url = short_url[:10]
	}

	return short_url
}

// generating a random alphanumeric symbol (including underscore) for collision handling
func GenerateRandomSymbol() string {
	const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	return string(symbols[rand.Intn(len(symbols))])
}