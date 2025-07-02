module github.com/structx/tbd/idp

go 1.24.4

replace (
	github.com/structx/tbd/lib/gateway => ../lib/go-gateway
	github.com/structx/tbd/lib/logging => ../lib/go-logging
	github.com/structx/tbd/lib/protocol => ../lib/go-protocol
	github.com/structx/tbd/lib/setup => ../lib/go-setup
	github.com/structx/tbd/lib/wallet => ../lib/go-wallet
)

require (
	buf.build/go/protovalidate v0.13.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/stretchr/testify v1.10.0
	github.com/structx/tbd/lib/gateway v0.0.0-00010101000000-000000000000
	github.com/structx/tbd/lib/logging v0.0.0-00010101000000-000000000000
	github.com/structx/tbd/lib/protocol v0.0.0-00010101000000-000000000000
	github.com/structx/tbd/lib/setup v0.0.0-00010101000000-000000000000
	github.com/structx/tbd/lib/wallet v0.0.0-00010101000000-000000000000
	go.uber.org/fx v1.24.0
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250613105001-9f2d3c737feb.1 // indirect
	cel.dev/expr v0.23.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/cel-go v0.25.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	go.dedis.ch/kyber/v4 v4.0.0-pre2 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
