package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

func GenerateShortUrl(original_url string) string {
	hash := sha1.Sum([]byte(original_url))

	short_url := base64.URLEncoding.EncodeToString(hash[:])
	short_url = strings.TrimRight(short_url, "=")

	if len(short_url) > 10 {
		short_url = short_url[:10]
	}

	return short_url
}

// collision handling
func GenerateRandomSymbol() string {
	const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	rand.Seed(time.Now().UnixNano())
	return string(symbols[rand.Intn(len(symbols))])
}