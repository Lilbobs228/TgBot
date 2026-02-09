package validators

import (
	"errors"
	"strings"
)

// ValidateWord проверяет валидность слова
func ValidateWord(word, definition, translation string) error {
	if strings.TrimSpace(word) == "" {
		return errors.New("слово не может быть пустым")
	}

	if strings.TrimSpace(definition) == "" {
		return errors.New("определение не может быть пустым")
	}

	if strings.TrimSpace(translation) == "" {
		return errors.New("перевод не может быть пустым")
	}

	// Проверяем длину слова (3-50 символов)
	if len(word) < 3 || len(word) > 50 {
		return errors.New("слово должно быть от 3 до 50 символов")
	}

	// Проверяем длину определения
	if len(definition) > 500 {
		return errors.New("определение слишком длинное (максимум 500 символов)")
	}

	// Проверяем длину перевода
	if len(translation) > 500 {
		return errors.New("перевод слишком длинный (максимум 500 символов)")
	}

	if containsDangerousChars(word) {
		return errors.New("слово содержит недопустимые символы")
	}

	if containsDangerousChars(definition) {
		return errors.New("определение содержит недопустимые символы")
	}

	if containsDangerousChars(translation) {
		return errors.New("перевод содержит недопустимые символы")
	}

	return nil
}

// ValidateBotToken проверяет валидность токена бота
func ValidateBotToken(token string) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("токен бота не может быть пустым")
	}

	// Проверяем формат токена Telegram бота (должен содержать двоеточие)
	if !strings.Contains(token, ":") {
		return errors.New("неверный формат токена бота")
	}

	// Проверяем длину токена
	if len(token) < 20 || len(token) > 100 {
		return errors.New("неверная длина токена бота")
	}

	return nil
}

// ValidateChatID проверяет валидность ID чата
func ValidateChatID(chatID int64) error {
	if chatID == 0 {
		return errors.New("ID чата не может быть равен нулю")
	}

	// Проверяем, что ID чата в разумных пределах для Telegram
	// Telegram ID чатов могут быть от -2^63 до 2^63-1, но ограничиваем разумными пределами
	if chatID < -999999999999999 || chatID > 999999999999999 {
		return errors.New("ID чата вне допустимого диапазона")
	}

	return nil
}

// containsDangerousChars проверяет наличие потенциально опасных символов
func containsDangerousChars(text string) bool {
	// Список потенциально опасных символов
	dangerousChars := []string{
		"<script", "</script>", "javascript:", "data:",
		"vbscript:", "onload=", "onerror=", "onclick=",
	}

	textLower := strings.ToLower(text)
	for _, char := range dangerousChars {
		if strings.Contains(textLower, char) {
			return true
		}
	}

	return false
}
