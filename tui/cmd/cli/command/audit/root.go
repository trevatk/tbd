package audit

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command"
)

var auditCmd = &cobra.Command{
	Use: "audit",
}

func init() {
	auditCmd.AddCommand(listCmd)
	command.RootCmd.AddCommand(auditCmd)
}
