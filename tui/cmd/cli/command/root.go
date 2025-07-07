package command

import (
	"context"
	"log/slog"
	"time"

	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/internal/pkg/logging"
)

const (
	defaultTimeout = 15
)

var (
	configFile string

	// RootCmd base cobra cli command
	RootCmd = &cobra.Command{
		Use: "zero",
	}
)

func init() {
	RootCmd.Flags().StringVarP(&configFile, "config-file", "c", "$HOME/.zerodot", "config file")
}

// Execute root cmd
func Execute(ctx context.Context) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*defaultTimeout)
	defer cancel()

	if err := RootCmd.ExecuteContext(timeout); err != nil {
		logging.FromContext(ctx).ErrorContext(ctx, "failed to execute command", slog.String("error", err.Error()))
	}
}
