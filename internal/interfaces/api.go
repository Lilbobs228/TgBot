package interfaces

import (
	"context"
	"telegram-quotes-bot/internal/entities"
)

// WordAPI определяет методы для работы с API слов
type WordAPI interface {
	GetRandomWord(ctx context.Context) (*entities.Word, error)
}

// TelegramSender определяет методы для отправки сообщений в Telegram.
type TelegramSender interface {
	SendMessage(ctx context.Context, message string) error
}

// CronScheduler определяет методы для планирования задач.
type CronScheduler interface {
	Start()
	AddJob(spec string, job func())
}
