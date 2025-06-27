package server

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command"
)

var (
	serverCmd = &cobra.Command{
		Use: "server",
	}
)

func init() {
	serverCmd.AddCommand(pingCmd)
	command.RootCmd.AddCommand(serverCmd)
}
