package thread

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command/chat"
)

var (
	threadCmd = &cobra.Command{
		Use: "thread",
	}
)

func init() {
	threadCmd.AddCommand(createCmd)
	chat.ChatCmd.AddCommand(threadCmd)
}
