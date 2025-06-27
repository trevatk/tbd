package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/structx/tbd/tui/cmd/cli/command"
	_ "github.com/structx/tbd/tui/cmd/cli/command/audit"
	_ "github.com/structx/tbd/tui/cmd/cli/command/chat"
	_ "github.com/structx/tbd/tui/cmd/cli/command/chat/thread"
	_ "github.com/structx/tbd/tui/cmd/cli/command/realm"
	_ "github.com/structx/tbd/tui/cmd/cli/command/server"
	_ "github.com/structx/tbd/tui/cmd/cli/command/user"
	"github.com/structx/tbd/tui/internal/pkg/logging"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := logging.New()
	ctx = logging.WithLogger(ctx, logger)

	command.Execute(ctx)
}
