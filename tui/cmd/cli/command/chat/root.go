package chat

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command"
)

var (
	ChatCmd = &cobra.Command{
		Use: "chat",
	}
)

func init() {
	command.RootCmd.AddCommand(ChatCmd)
}
