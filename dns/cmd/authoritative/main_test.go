package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"soft.structx.io/idp/lib/logging"
	"soft.structx.io/idp/lib/setup"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/structx/tbd/lib/protocol"
	"github.com/trevatk/tbd/dns/internal/authoritative"

	pb "github.com/structx/tbd/lib/protocol/dns/kademlia/v1"
)

func TestAuthoritativeMain(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	cfg := setup.UnmarshalConfig()
	logger := logging.New(cfg.Logger.Level)

	kv := authoritative.NewKv()
	dht := authoritative.NewDHT(kv, cfg.Gateway.Host, cfg.Gateway.Port)
	trs := authoritative.NewTransport(logger, dht)

	opts := []protocol.TestServerOption{
		protocol.WithTestTransports(trs),
		protocol.WithTestLogger(logger),
	}

	ts := protocol.NewTestServer(opts...)
	go ts.Start(ctx)
	defer ts.Stop(ctx)

	conn, err := protocol.NewTestConn(ctx, ts.BufDialer)
	if err != nil {
		t.Fatalf("failed to create test conn: %v", err)
	}
	client := pb.NewKademliaServiceClient(conn)

	runIntegrationTests(t, ctx, client)
}

func runIntegrationTests(t *testing.T, ctx context.Context, client pb.KademliaServiceClient) {
	assert := assert.New(t)

	var (
		expected error
	)

	expected = nil
	err := pingNode(ctx, client)
	assert.Equal(expected, err)
}

func pingNode(ctx context.Context, client pb.KademliaServiceClient) error {
	_, err := client.Ping(ctx, &pb.PingRequest{
		Sender: &pb.Node{
			NodeId:     "4b84b15bff6ee5796152495a230e45e3d7e947d9",
			IpOrDomain: "127.0.0.1",
			Port:       53,
			LastSeen:   timestamppb.Now(),
		},
		RequestId: uuid.New().String(),
	})
	if err != nil {
		return fmt.Errorf("failed to execute ping command: %w", err)
	}
	return nil
}
