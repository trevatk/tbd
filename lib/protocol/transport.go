package protocol

import "google.golang.org/grpc"

// Transport protocol transport params
type Transport struct {
	ServiceDesc *grpc.ServiceDesc
	Service     any
}
