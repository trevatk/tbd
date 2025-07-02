// Package main entrypoint for dapp server
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.structx/tbd/dapp/internal"

	"github.com/structx/tbd/lib/logging"
	"github.com/structx/tbd/lib/protocol"
	"github.com/structx/tbd/lib/setup"
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
