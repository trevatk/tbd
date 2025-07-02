package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/trevatk/tbd/dns/internal/resolver"

	"github.com/structx/tbd/lib/logging"
	"github.com/structx/tbd/lib/protocol"
	"github.com/structx/tbd/lib/setup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGINT)
	defer func() {
		cancel()

		if r := recover(); r != nil {
			slog.Default().Error("panic recovery", r)
		}
	}()

	if err := realMain(ctx); err != nil {
		slog.ErrorContext(ctx, "start resolver", "error", err)
	}
}

func realMain(ctx context.Context) error {
	cfg := setup.UnmarshalConfig()
	logger := logging.New(cfg.Logger.Level)

	cache := resolver.NewCache()
	cache.Start()
	defer cache.Stop()

	ns := []string{cfg.Nameserver.NS1, cfg.Nameserver.NS2}
	tr := resolver.NewTransport(logger, ns, cache)
	trs := []protocol.Transport{tr}

	opts := []protocol.ServerOption{
		protocol.WithHost(cfg.Gateway.Host),
		protocol.WithPort(cfg.Gateway.Port),
		protocol.WithTransports(trs),
		protocol.WithLogger(logger),
	}

	s := protocol.NewServer(opts...)
	return s.StartAndStop(ctx)
}
