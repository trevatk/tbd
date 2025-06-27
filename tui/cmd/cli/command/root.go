package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	configFile string

	RootCmd = &cobra.Command{
		Use: "zero",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	RootCmd.Flags().StringVarP(&configFile, "config-file", "c", "$HOME/.zerodot", "config file")
}

func Execute(ctx context.Context) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	if err := RootCmd.ExecuteContext(timeout); err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute cli %v\n", err)
		os.Exit(1)
	}
}
