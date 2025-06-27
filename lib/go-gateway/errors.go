package gateway

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ErrInvalidArgument ...
func ErrInvalidArgument() error {
	return status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
}

// ErrInternal ...
func ErrInternal() error {
	return status.Error(codes.Internal, codes.Internal.String())
}

// NewEmptyResponse ...
func NewEmptyResponse() *emptypb.Empty {
	return &emptypb.Empty{}
}
