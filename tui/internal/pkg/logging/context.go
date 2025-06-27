package logging

import (
	"context"
	"log/slog"
)

type contextKey string

const loggerKey contextKey = "logger"

// WithLogger attach logger to context
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext extract logger from context
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return New()
}
