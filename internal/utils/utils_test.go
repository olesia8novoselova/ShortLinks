package utils

import (
	"testing"
	"unicode"
)

func TestGenerateShortUrl(t *testing.T) {
	// тестовые случаи
	tests := []struct {
		name        string
		originalUrl string
	}{
		{"Простой URL", "https://example.com"},
		{"Длинный URL", "https://www.google.com/search?q=this+is+a+very+long+url+with+many+parameters"},
		{"Короткий URL", "https://short.url"},
		{"URL со специальными символами", "https://example.com/path?query=param&value=123!@#"},
		{"Пустой URL", ""}, // пустой ввод
	}

	// map для отслеживания уникальности коротких URL
	shortUrlMap := make(map[string]string)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			short := GenerateShortUrl(tt.originalUrl)

			// проверка длины
			if len(short) != 10 {
				t.Errorf("Тестовый случай '%s': Ожидалась длина короткого URL 10, получено %d", tt.name, len(short))
			}

			// проверка набора символов
			for _, char := range short {
				if !unicode.IsLower(char) && !unicode.IsUpper(char) && !unicode.IsDigit(char) && char != '_' {
					t.Errorf("Тестовый случай '%s': Недопустимый символ '%c' в коротком URL. Разрешены только a-z, A-Z, 0-9 и _.", tt.name, char)
				}
			}

			// проверка уникальности
			if existingOriginal, exists := shortUrlMap[short]; exists {
				t.Errorf("Тестовый случай '%s': Дублирующийся короткий URL '%s' сгенерирован для '%s' и '%s'", tt.name, short, existingOriginal, tt.originalUrl)
			} else {
				shortUrlMap[short] = tt.originalUrl
			}

			// проверка детерминированного поведения
			short2 := GenerateShortUrl(tt.originalUrl)
			if short != short2 {
				t.Errorf("Тестовый случай '%s': Не детерминированное поведение. Сгенерированы '%s' и '%s' для одного и того же ввода.", tt.name, short, short2)
			}
		})
	}
}
