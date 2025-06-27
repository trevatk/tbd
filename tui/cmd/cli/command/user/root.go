package user

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/cmd/cli/command"
)

var (
	userCmd = &cobra.Command{
		Use: "user",
	}
)

func init() {
	createCmd.Flags().StringVarP(&userEmail, "email", "e", "", "user email")
	userCmd.AddCommand(createCmd)
	command.RootCmd.AddCommand(userCmd)
}
