package protocol

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

// NewConn return new gRPC client connection
func NewConn(target string) (*grpc.ClientConn, error) {
	return grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// NewConn return new gRPC client connection
func NewConnWithResolver(target string, builder resolver.Builder) (*grpc.ClientConn, error) {
	return grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithResolvers(builder))
}
