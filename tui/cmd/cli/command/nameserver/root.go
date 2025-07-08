package nameserver

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command"
)

var (
	Cmd = &cobra.Command{
		Use:     "nameserver",
		Aliases: []string{"ns"},
	}
)

func init() {
	command.RootCmd.AddCommand(Cmd)
}
