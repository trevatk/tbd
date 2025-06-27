package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/structx/tbd/lib/gateway"
	"github.com/structx/tbd/lib/logging"
	"github.com/structx/tbd/lib/setup"

	"github.com/structx/tbd/idp/internal/audit"
	"github.com/structx/tbd/idp/internal/audit/lsm"
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

	s := gateway.New(opts...)

	return s.StartAndStop(ctx)
}
