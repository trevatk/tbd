package record

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command/nameserver"
)

var (
	cmd = &cobra.Command{
		Use: "record",
	}
)

func init() {
	cmd.AddCommand(createCmd)
	cmd.AddCommand(listCmd)
	nameserver.Cmd.AddCommand(cmd)
}
