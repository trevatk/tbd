package thread

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/structx/tbd/tui/internal/pkg/logging"
	"github.com/trevatk/tbd/lib/protocol"
	pb "github.com/trevatk/tbd/lib/protocol/chat/v1"
)

const (
	defaultTimeout = 250
)

var (
	createCmd = &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			client, err := newClient(serverAddr)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			timeout, cancel := context.WithTimeout(ctx, time.Millisecond*defaultTimeout)
			defer cancel()

			resp, err := client.CreateThread(timeout, &pb.CreateThreadRequest{
				DisplayName: "hello",
			})
			if err != nil {
				return fmt.Errorf("failed to create thread: %w", err)
			}

			logging.FromContext(ctx).Info("thread successfully created...", "thread_id", resp.Id)

			return nil
		},
	}
)

func newClient(target string) (pb.ChatServiceClient, error) {
	conn, err := protocol.NewConn(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}
	return pb.NewChatServiceClient(conn), nil
}
