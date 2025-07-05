package thread

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/internal/pkg/logging"
	"github.com/trevatk/tbd/lib/protocol/chat"
)

const (
	defaultTimeout = 250
)

var (
	createCmd = &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			client, err := chat.NewClient(serverAddr)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			timeout, cancel := context.WithTimeout(ctx, time.Millisecond*defaultTimeout)
			defer cancel()

			threadID, err := client.CreateThread(timeout, "")
			if err != nil {
				return fmt.Errorf("failed to create thread: %w", err)
			}

			logging.FromContext(ctx).Info("thread successfully created...", "thread_id", threadID)

			return nil
		},
	}
)
