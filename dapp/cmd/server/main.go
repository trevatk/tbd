package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/structx/tbd/lib/gateway"
	"github.com/structx/tbd/lib/logging"
	"github.com/structx/tbd/lib/setup"
	"github.structx/tbd/dapp/internal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		cancel()

		if r := recover(); r != nil {
			log.Fatalf("panic recovery %v", r)
		}
	}()

	if err := realMain(ctx); err != nil {
		log.Fatal(err)
	}
}

func realMain(ctx context.Context) error {
	cfg := setup.UnmarshalConfig()
	logger := logging.New(cfg.Logger.Level)
	desc, service := internal.NewTransport(logger)

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
