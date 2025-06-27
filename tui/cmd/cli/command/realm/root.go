package realm

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command"
)

var (
	realmCmd = &cobra.Command{
		Use: "realm",
	}
)

func init() {
	realmCmd.AddCommand(createCmd)
	command.RootCmd.AddCommand(realmCmd)
}
