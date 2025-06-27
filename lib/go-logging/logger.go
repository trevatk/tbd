package logging

import (
	"log/slog"
	"os"
	"strings"
)

// New
func New(level string) *slog.Logger {

	var leveler slog.Leveler
	switch strings.ToLower(level) {
	case "info":
		leveler = slog.LevelInfo
	case "error":
		leveler = slog.LevelError
	case "warn":
		leveler = slog.LevelWarn
	default:
		leveler = slog.LevelDebug
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: leveler,
	}))
}