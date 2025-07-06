package wallet

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command"
)

var (
	walletCmd = &cobra.Command{
		Use: "wallet",
	}
)

func init() {
	walletCmd.AddCommand(createCmd)
	command.RootCmd.AddCommand(walletCmd)
}
