package protocol

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

// ErrAlreadyExists ...
func ErrAlreadyExists() error {
	return status.Error(codes.AlreadyExists, codes.AlreadyExists.String())
}

// ErrNotFound ...
func ErrNotFound() error {
	return status.Error(codes.NotFound, codes.NotFound.String())
}

// NewEmptyResponse ...
func NewEmptyResponse() *emptypb.Empty {
	return &emptypb.Empty{}
}
