package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	dapp "github.com/trevatk/tbd/dapp/internal"
	"github.com/trevatk/tbd/lib/logging"
	"github.com/trevatk/tbd/lib/protocol"
	"github.com/trevatk/tbd/lib/setup"

	pb "github.com/trevatk/tbd/lib/protocol/chat/v1"
)

func TestDappMain(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second*15)
	defer cancel()

	cfg := setup.UnmarshalConfig()
	logger := logging.New(cfg.Logger.Level)

	svc := dapp.NewService()
	tr := dapp.NewTransport(logger, svc)

	opts := []protocol.TestServerOption{
		protocol.WithTestLogger(logger),
		protocol.WithTestTransports([]protocol.Transport{tr}),
	}

	ts := protocol.NewTestServer(opts...)
	go ts.Start(ctx)
	defer ts.Stop(ctx)

	conn, err := protocol.NewTestConn(ctx, ts.BufDialer)
	if err != nil {
		t.Fatalf("failed to create test conn: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)
	runIntegrationTests(t, ctx, client)
}

func runIntegrationTests(t *testing.T, ctx context.Context, client pb.ChatServiceClient) {
	assert := assert.New(t)

	var (
		expected error
	)

	expected = nil
	err := createThread(ctx, client)
	assert.Equal(expected, err)
}

func createThread(ctx context.Context, client pb.ChatServiceClient) error {
	var (
		threadName = "helloworld"

		request = &pb.CreateThreadRequest{
			Members:     []string{},
			DisplayName: threadName,
		}
	)

	resp, err := client.CreateThread(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to execute create thread gRPC call: %w", err)
	}

	if resp.Thread.DisplayName != threadName {
		return fmt.Errorf("unexpected thread name %s expected %s", resp.Thread.DisplayName, threadName)
	}

	_, err = uuid.Parse(resp.Thread.Id)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", err)
	}

	return nil
}
