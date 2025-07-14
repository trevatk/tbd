package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/trevatk/tbd/dns/internal/nameserver"
	"github.com/trevatk/tbd/lib/logging"
	"github.com/trevatk/tbd/lib/protocol"
	"github.com/trevatk/tbd/lib/setup"
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

	kv := nameserver.NewKv()
	dht := nameserver.NewDHT(kv, cfg.Gateway.Host, cfg.Gateway.Port)
	trs := nameserver.NewTransport(logger, dht)

	as, err := nameserver.NewAuthoritativeServer(logger, dht)
	if err != nil {
		return fmt.Errorf("failed to create authoritative server: %w", err)
	}

	if err = nameserver.AddRecordsFromJson(cfg.Records.JsonPath, logger, dht); err != nil {
		return fmt.Errorf("failed to add records from json: %w", err)
	}

	opts := []protocol.ServerOption{
		protocol.WithHost(cfg.Gateway.Host),
		protocol.WithPort(cfg.Gateway.Port),
		protocol.WithTransports(trs),
		protocol.WithLogger(logger),
	}

	_ = protocol.NewServer(opts...)

	// start udp server
	return as.Listen(ctx)

	// start gRPC server
	// return s.StartAndStop(ctx)
}
