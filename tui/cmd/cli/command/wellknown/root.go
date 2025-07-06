package wellknown

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command"
)

var (
	wellknownCmd = &cobra.Command{
		Use:     "wellknown",
		Aliases: []string{"w"},
		Short:   "wellknown configuration tool",
	}
)

func init() {
	wellknownCmd.AddCommand(generateCmd)
	command.RootCmd.AddCommand(wellknownCmd)
}
