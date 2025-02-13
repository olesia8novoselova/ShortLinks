package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// допустимый набор символов для короткого URL
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func GenerateShortUrl(originalUrl string) string {
	// вычисляем SHA256 хэш оригинального URL
	hash := sha256.Sum256([]byte(originalUrl))
	// преобразуем хэш в шестнадцатеричную строку
	hashHex := hex.EncodeToString(hash[:])
	// преобразуем шестнадцатеричную строку в большое целое число
	hashInt := new(big.Int)
	hashInt.SetString(hashHex, 16)

	var shortUrl string
	for i := 0; i < 10; i++ {
		// вычисляем остаток от деления на длину набора символов
		remainder := new(big.Int).Mod(hashInt, big.NewInt(int64(len(charset))))
		// добавляем символ из набора charset в shortUrl
		shortUrl += string(charset[remainder.Int64()])
		// делим hashInt на длину набора символов для следующей итерации
		hashInt.Div(hashInt, big.NewInt(int64(len(charset))))
	}
	return shortUrl
}
