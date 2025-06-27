package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/structx/tbd/dns/internal/resolver"
	"github.com/structx/tbd/lib/gateway"
	"github.com/structx/tbd/lib/logging"
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
	trs := []gateway.Transport{tr}

	opts := []gateway.Option{
		gateway.WithHost(cfg.Gateway.Host),
		gateway.WithPort(cfg.Gateway.Port),
		gateway.WithTransports(trs),
		gateway.WithLogger(logger),
	}

	s := gateway.New(opts...)
	return s.StartAndStop(ctx)
}
