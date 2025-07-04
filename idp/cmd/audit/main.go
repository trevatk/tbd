package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/trevatk/tbd/lib/logging"
	"github.com/trevatk/tbd/lib/protocol"
	"github.com/trevatk/tbd/lib/setup"

	"github.com/trevatk/tbd/idp/internal/audit"
	"github.com/trevatk/tbd/idp/internal/audit/lsm"
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

	lsm, err := lsm.New(cfg.KeyValue.Dir)
	if err != nil {
		return fmt.Errorf("failed to initialize lsm: %w", err)
	}

	svc, err := audit.NewService(nil, lsm)
	if err != nil {
		return fmt.Errorf("failed to initialize audit service: %w", err)
	}

	desc, service := audit.NewTransport(logger, svc)

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
