package protocol

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

// NewTestConn create new test connection with bufnet dialer
// and passthrough resolver scheme
func NewTestConn(
	ctx context.Context,
	dialer func(context.Context, string) (net.Conn, error),
) (*grpc.ClientConn, error) {
	resolver.SetDefaultScheme("passthrough")
	return grpc.NewClient(
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
