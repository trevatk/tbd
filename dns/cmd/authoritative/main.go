package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/trevatk/tbd/dns/internal/authoritative"

	"github.com/structx/tbd/lib/gateway"
	"github.com/structx/tbd/lib/logging"
	"github.com/structx/tbd/lib/setup"
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

	kv := authoritative.NewKv()
	dht := authoritative.NewDHT(kv, cfg.Gateway.Host, cfg.Gateway.Port)
	trs := authoritative.NewTransport(logger, dht)

	opts := []gateway.Option{
		gateway.WithHost(cfg.Gateway.Host),
		gateway.WithPort(cfg.Gateway.Port),
		gateway.WithTransports(trs),
		gateway.WithLogger(logger),
	}

	s := gateway.New(opts...)
	return s.StartAndStop(ctx)
}
