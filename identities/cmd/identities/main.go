package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/structx/tbd/identities/internal/identities"
	"github.com/structx/tbd/lib/gateway"
	"github.com/structx/tbd/lib/logging"
	"github.com/structx/tbd/lib/setup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		cancel()

		if r := recover(); r != nil {
			log.Fatalf("panic recover: %v", r)
		}
	}()

	if err := realMain(ctx); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func realMain(ctx context.Context) error {
	cfg := setup.UnmarshalConfig()

	logger := logging.New(cfg.Logger.Level)
	logger.InfoContext(ctx, "service configuration", slog.Any("config", cfg))
	_ = identities.NewAuth(cfg.Auth.SigningKey)
	graph := identities.NewGraph()
	svc := identities.NewService(graph)
	desc, service := identities.NewTransport(logger, svc)

	trs := []gateway.Transport{
		{
			ServiceDesc: desc,
			Service:     service,
		},
	}

	opts := []gateway.Option{
		gateway.WithHost(cfg.Gateway.Host),
		gateway.WithPort(cfg.Gateway.Port),
		gateway.WithTransports(trs),
		gateway.WithLogger(logger),
	}

	server := gateway.New(opts...)

	return server.StartAndStop(ctx)
}
