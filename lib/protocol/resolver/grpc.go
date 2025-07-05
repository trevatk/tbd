package resolver

import (
	"log/slog"

	"google.golang.org/grpc/resolver"
)

type builder struct {
	scheme      string
	serviceName string

	logger *slog.Logger
}

// Option
type Option func(*builder)

// Register
// func Register(builder *builder) {
// 	resolver.Register(builder)
// }

// NewBuilder returns new gRPC resolver builder
func NewBuilder(opts ...Option) *builder {
	builder := &builder{}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

// WithLogger
func WithLogger(logger *slog.Logger) Option {
	return func(b *builder) {
		b.logger = logger
	}
}

// WithScheme
func WithScheme(scheme string) Option {
	return func(b *builder) {
		b.scheme = scheme
	}
}

// WithServiceName
func WithServiceName(name string) Option {
	return func(b *builder) {
		b.serviceName = name
	}
}

// Build gRPC resolver
func (b *builder) Build(
	target resolver.Target,
	oc resolver.ClientConn,
	_ resolver.BuildOptions,
) (resolver.Resolver, error) {
	r := &nameResolver{}
	r.start()
	return r, nil
}

// Scheme getter scheme
func (b *builder) Scheme() string {
	return b.scheme
}

type nameResolver struct {
	target    resolver.Target
	cc        resolver.ClientConn
	addrStore map[string][]string
}

func (r *nameResolver) start() {
	addrStrs := r.addrStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, addr := range addrStrs {
		addrs[i] = resolver.Address{Addr: addr}
	}
	_ = r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (r *nameResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *nameResolver) Close()                                {}
