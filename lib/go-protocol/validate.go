package protocol

import (
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

// Validate proto message
func Validate(m proto.Message) error {
	return protovalidate.Validate(m, protovalidate.WithFailFast())
}
