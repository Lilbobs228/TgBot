package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"telegram-quotes-bot/internal/adapters"
	"telegram-quotes-bot/internal/config"
	"telegram-quotes-bot/internal/usecases"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

// setupLogger логгер
func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger
}

func main() {
	// Настройка логгера
	logger := setupLogger()

	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		logger.Warn("Файл .env не найден или не загружен")
	}

	// Создаем контекст с обработкой сигналов для graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Error("Ошибка загрузки конфигурации", "error", err)
		os.Exit(1)
	}

	// Инициализация адаптеров
	// quoteAPI := adapters.NewForismaticAPI()
	wordAPI := adapters.NewDictionaryAPI()
	telegramAdapter, err := adapters.NewTelegramAdapter(cfg.BotToken, cfg.ChatID)
	if err != nil {
		logger.Error("Не удалось инициализировать TelegramAdapter", "error", err)
		os.Exit(1)
	}

	// Инициализация сервисов
	// fetchQuoteService := usecases.NewFetchQuoteService(quoteAPI)
	// sendQuoteService := usecases.NewSendQuoteService(telegramAdapter)
	fetchWordService := usecases.NewFetchWordService(wordAPI)
	sendWordService := usecases.NewSendWordService(telegramAdapter)

	// Планировщик Cron
	c := cron.New()
	defer c.Stop()

	// Задача отправки слов
	_, err = c.AddFunc("0 4,8,14,18 * * *", func() {
		taskCtx := context.Background()

		// Получение слова
		word, err := fetchWordService.FetchWord(taskCtx)
		if err != nil {
			logger.Error("Ошибка получения слова", "error", err)
			return
		}

		// Отправка слова
		if err := sendWordService.SendWord(taskCtx, word); err != nil {
			logger.Error("Ошибка отправки слова", "error", err)
		} else {
			logger.Info("Слово успешно отправлено", "word", word.Word, "definition", word.Definition)
		}
	})
	if err != nil {
		logger.Error("Не удалось добавить cron-задачу", "error", err)
		os.Exit(1)
	}

	// Запуск планировщика
	c.Start()
	logger.Info("Планировщик запущен. Ожидание задач.")

	// Отправка тестового слова при запуске (если включено в конфигурации)
	if cfg.SendTestQuote {
		logger.Info("Отправка тестового слова...")
		testCtx := context.Background()

		// Получение тестового слова
		testWord, err := fetchWordService.FetchWord(testCtx)
		if err != nil {
			logger.Error("Ошибка получения тестового слова", "error", err)
		} else {
			// Отправка тестового слова
			if err := sendWordService.SendWord(testCtx, testWord); err != nil {
				logger.Error("Ошибка отправки тестового слова", "error", err)
			} else {
				logger.Info("Тестовое слово успешно отправлено", "word", testWord.Word, "definition", testWord.Definition)
			}
		}
	} else {
		logger.Info("Отправка тестового слова отключена в конфигурации")
	}

	// Ожидание сигнала завершения
	<-ctx.Done()
	logger.Info("Получен сигнал завершения. Останавливаем планировщик...")

	// Остановка планировщика
	stopCtx := c.Stop()
	<-stopCtx.Done()
	logger.Info("Планировщик остановлен. Программа завершена.")
}
