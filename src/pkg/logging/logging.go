package logging

import (
	"log/slog"
	"os"
)

// Определяем типы окружений для логгера
const (
	envLocal = "local"
	envProd  = "prod"
)

// Logger структура для работы с логированием
type Logger struct {
	log *slog.Logger
}

// NewLogger создает новый логгер в зависимости от окружения (local или prod)
func NewLogger(env string) *Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		// Логгирование в текстовом формате для локальной среды с уровнем debug
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		// Логгирование в JSON формате для продакшена с уровнем info
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		// По умолчанию используем текстовый формат с уровнем info
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return &Logger{log: logger}
}

// Info логгирует информационное сообщение
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Info(msg, keysAndValues...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log.Error(msg, keysAndValues...)
	os.Exit(1) // Завершение программы с кодом 1
}

// Error логгирует ошибку
func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.log.Error(msg, keysAndValues...)
}

// Debug логгирует отладочные сообщения (доступно только при уровне debug)
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.log.Debug(msg, keysAndValues...)
}

// Warn логгирует предупреждающие сообщения
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.log.Warn(msg, keysAndValues...)
}

// WithContext добавляет контекст к логам (например, request_id)
func (l *Logger) WithContext(key string, value interface{}) *Logger {
	return &Logger{log: l.log.With(key, value)}
}

// With добавляет дополнительные поля к логам
func (l *Logger) With(key string, value interface{}) *Logger {
	return &Logger{log: l.log.With(key, value)}
}
