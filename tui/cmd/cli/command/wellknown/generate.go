package wellknown

import (
	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/internal/pkg/logging"
)

var (
	generateCmd = &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logging.FromContext(ctx).InfoContext(ctx, "wellknown generate handler")
			return nil
		},
	}
)
