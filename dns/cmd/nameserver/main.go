package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/trevatk/tbd/dns/internal/nameserver"
	"github.com/trevatk/tbd/lib/logging"
	"github.com/trevatk/tbd/lib/protocol"
	"github.com/trevatk/tbd/lib/setup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer func() {
		cancel()

		if r := recover(); r != nil {
			slog.Error("panic recover: %v", r)
		}
	}()

	if err := realMain(ctx); err != nil {
		slog.Error("application run failed", "error", err)
	}
}

func realMain(ctx context.Context) error {
	cfg := setup.UnmarshalConfig()
	logger := logging.New(cfg.Logger.Level)

	kv := nameserver.NewKv()
	dht := nameserver.NewDHT(kv, cfg.Gateway.Host, cfg.Gateway.Port)
	trs := nameserver.NewTransport(logger, dht)

	opts := []protocol.ServerOption{
		protocol.WithHost(cfg.Gateway.Host),
		protocol.WithPort(cfg.Gateway.Port),
		protocol.WithTransports(trs),
		protocol.WithLogger(logger),
	}

	s := protocol.NewServer(opts...)
	return s.StartAndStop(ctx)
}
