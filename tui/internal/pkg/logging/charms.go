package logging

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

// New
func New() *slog.Logger {
	handler := log.New(os.Stderr)
	return slog.New(handler)
}
