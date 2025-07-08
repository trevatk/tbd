package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/trevatk/tbd/lib/logging"
	"github.com/trevatk/tbd/lib/protocol"
	"github.com/trevatk/tbd/lib/setup"
	wellknown "github.com/trevatk/tbd/wellknown/internal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		cancel()

		if r := recover(); r != nil {
			slog.Default().Error("panic recovery", slog.Any("panic", r))
		}
	}()

	if err := realMain(ctx); err != nil {
		slog.Default().ErrorContext(ctx, "application failture", slog.String("error", err.Error()))
	}
}

func realMain(ctx context.Context) error {
	cfg := setup.UnmarshalConfig()
	logger := logging.New(cfg.Logger.Level)

	tr := wellknown.NewTransport(logger, cfg.Wellknown.Path)
	trs := []protocol.Transport{tr}

	opts := []protocol.ServerOption{
		protocol.WithHost(cfg.Gateway.Host),
		protocol.WithPort(cfg.Gateway.Port),
		protocol.WithLogger(logger),
		protocol.WithTransports(trs),
	}

	s := protocol.NewServer(opts...)

	return s.StartAndStop(ctx)
}
