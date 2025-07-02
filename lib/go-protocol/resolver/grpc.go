package resolver

import "google.golang.org/grpc/resolver"

// Builder gRPC resolver builder
type Builder struct {
	scheme      string
	serviceName string
}

// NewBuilder returns new gRPC resolver builder
func NewBuilder(scheme, serviceName string) *Builder {
	builder := &Builder{
		scheme:      scheme,
		serviceName: serviceName,
	}
	resolver.Register(builder)
	return builder
}

// Build gRPC resolver
func (b *Builder) Build(
	target resolver.Target,
	oc resolver.ClientConn,
	_ resolver.BuildOptions,
) (resolver.Resolver, error) {
	r := &nameResolver{}
	r.start()
	return r, nil
}

// Scheme getter scheme
func (b *Builder) Scheme() string {
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
