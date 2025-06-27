package resolver

import (
	"fmt"
	"io"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/structx/tbd/lib/protocol/dns/resolver/v1"
)

func resolveDID(target string) ([]byte, error) {

	resp, err := http.Get(target)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http get: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode > http.StatusAccepted {
		return nil, fmt.Errorf("unexcepted http response %d %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func newResponseWithError(code codes.Code) *pb.ResolveResponse {
	var status pb.ResolveResponse_ResponseStatus
	switch code {
	case codes.NotFound:
		status = pb.ResolveResponse_RESPONSE_STATUS_DID_NOT_FOUND
	}

	return &pb.ResolveResponse{
		Status:       status,
		ErrorMessage: code.String(),
	}
}

func newAuthoritativeClient(target string) (pb.DNSResolverServiceClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return pb.NewDNSResolverServiceClient(conn), nil
}
