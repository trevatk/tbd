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

	output      string
	alsoKnownAs []string
)

func init() {
	generateCmd.PersistentFlags().StringVarP(&output, "output", "o", ".", "set output")
	generateCmd.PersistentFlags().StringArrayVarP(&alsoKnownAs, "also known as", "a", []string{}, "set also known as")

	wellknownCmd.AddCommand(generateCmd)
	command.RootCmd.AddCommand(wellknownCmd)
}
