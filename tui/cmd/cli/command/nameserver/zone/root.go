package zone

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command/nameserver"
)

var (
	cmd = &cobra.Command{
		Use: "zone",
	}
)

func init() {
	cmd.AddCommand(createCmd)
	nameserver.Cmd.AddCommand(cmd)
}
