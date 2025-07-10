package main

import (
	"context"
	"log/slog"

	dapp "github.com/trevatk/tbd/dapp/internal"
	"github.com/trevatk/tbd/lib/logging"
	"github.com/trevatk/tbd/lib/protocol"
	"github.com/trevatk/tbd/lib/setup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()

		if r := recover(); r != nil {
			slog.Error("panic recovery", slog.Any("panic", r))
		}
	}()

	if err := realMain(ctx); err != nil {
		slog.ErrorContext(ctx, "application error", slog.String("error", err.Error()))
	}
}

func realMain(ctx context.Context) error {
	cfg := setup.UnmarshalConfig()
	logger := logging.New(cfg.Logger.Level)

	svc := dapp.NewService()
	tr := dapp.NewTransport(logger, svc)

	opts := []protocol.ServerOption{
		protocol.WithHost(cfg.Gateway.Host),
		protocol.WithPort(cfg.Gateway.Port),
		protocol.WithLogger(logger),
		protocol.WithTransports([]protocol.Transport{tr}),
	}
	s := protocol.NewServer(opts...)

	return s.StartAndStop(ctx)
}
