package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	foundations "github.com/structx/tbd/dns/internal/authority"
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

	dht := foundations.NewDHT(cfg.Gateway.Host, cfg.Gateway.Port)
	trs := foundations.NewTransport(logger, dht)

	opts := []gateway.Option{
		gateway.WithHost(cfg.Gateway.Host),
		gateway.WithPort(cfg.Gateway.Port),
		gateway.WithTransports(trs),
		gateway.WithLogger(logger),
	}

	s := gateway.New(opts...)
	return s.StartAndStop(ctx)
}
