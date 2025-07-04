package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/trevatk/tbd/idp/internal/identities"
	"github.com/trevatk/tbd/lib/logging"
	"github.com/trevatk/tbd/lib/protocol"
	"github.com/trevatk/tbd/lib/setup"
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
}

func realMain(ctx context.Context) error {
	cfg := setup.UnmarshalConfig()

	logger := logging.New(cfg.Logger.Level)
	logger.InfoContext(ctx, "service configuration", slog.Any("config", cfg))
	_ = identities.NewAuth(cfg.Auth.SigningKey)
	graph := identities.NewGraph()
	svc := identities.NewService(graph)
	desc, service := identities.NewTransport(logger, svc)

	trs := []protocol.Transport{
		{
			ServiceDesc: desc,
			Service:     service,
		},
	}

	opts := []protocol.ServerOption{
		protocol.WithHost(cfg.Gateway.Host),
		protocol.WithPort(cfg.Gateway.Port),
		protocol.WithTransports(trs),
		protocol.WithLogger(logger),
	}

	s := protocol.NewServer(opts...)

	return s.StartAndStop(ctx)
}
