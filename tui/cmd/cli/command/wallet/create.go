package wallet

import (
	"github.com/spf13/cobra"
	"go.dedis.ch/kyber/v4/group/edwards25519"

	"github.com/trevatk/tbd/lib/wallet"
)

var (
	createCmd = &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := edwards25519.NewBlakeSHA256Ed25519WithRand(nil)
			w := wallet.NewV1(s)
			return w.Export("wallet.json")
		},
	}
)
