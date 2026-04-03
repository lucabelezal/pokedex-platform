package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Setup inicializa o logger estruturado global da aplicacao.
// Le LOG_LEVEL (debug/info/warn/error, padrao info) e LOG_FORMAT (json padrao, text para dev).
func Setup(service string) *slog.Logger {
	level := parseLevel(os.Getenv("LOG_LEVEL"))

	var handler slog.Handler
	if strings.EqualFold(strings.TrimSpace(os.Getenv("LOG_FORMAT")), "text") {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	logger := slog.New(handler).With("service", service)
	slog.SetDefault(logger)
	return logger
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
