package protocol

import (
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

// Validate
func Validate(m proto.Message) error {
	return protovalidate.Validate(m, protovalidate.WithFailFast())
}
