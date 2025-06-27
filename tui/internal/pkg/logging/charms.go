package logging

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

// New return new charms logger with slog
func New() *slog.Logger {
	handler := log.New(os.Stderr)
	return slog.New(handler)
}
