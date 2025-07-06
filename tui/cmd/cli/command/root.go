package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
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
		_, _ = fmt.Fprintf(os.Stderr, "failed to execute cli %v\n", err)
	}
}
