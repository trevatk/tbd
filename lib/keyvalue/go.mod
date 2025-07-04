module github.com/trevatk/tbd/lib/keyvalue

go 1.24.4

replace github.com/trevatk/tbd/lib/protocol => ../protocol

require (
	github.com/stretchr/testify v1.10.0
	github.com/trevatk/tbd/lib/protocol v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
