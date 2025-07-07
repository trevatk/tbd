module github.com/trevatk/tbd/wellknown

go 1.24.4

replace (
	github.com/trevatk/tbd/lib/logging => ../lib/logging
	github.com/trevatk/tbd/lib/protocol => ../lib/protocol
	github.com/trevatk/tbd/lib/setup => ../lib/setup
)

require (
	github.com/stretchr/testify v1.10.0
	github.com/trevatk/tbd/lib/logging v0.0.0-00010101000000-000000000000
	github.com/trevatk/tbd/lib/protocol v0.0.0-00010101000000-000000000000
	github.com/trevatk/tbd/lib/setup v0.0.0-00010101000000-000000000000
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250625184727-c923a0c2a132.1 // indirect
	buf.build/go/protovalidate v0.13.1 // indirect
	cel.dev/expr v0.23.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/google/cel-go v0.25.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/grpc v1.73.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
