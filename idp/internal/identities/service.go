package identities

import (
	"context"
	"fmt"
)

type whoami struct {
	P []byte // public key
}

type realmCreate struct {
	name string
	who  whoami
}

type realm struct {
	hash string
	name string
}

type userCreate struct {
	realm string
	email string
}

type user struct {
	hash string
}

type service struct {
	g *graph
}

// NewService return new access service implementation
func NewService(graph *graph) *service {
	return &service{
		g: graph,
	}
}

func (s *service) createRealm(_ context.Context, create realmCreate) (realm, error) {
	// TODO
	// verify who has permissions
	// to create realm

	_, _ = s.g.getVertex("")

	attrs := make(map[string]interface{})
	if err := s.g.addVertex(&vertex{
		id:         create.name,
		resource:   "REALM",
		edges:      []*edge{},
		attributes: attrs,
	}); err != nil {
		return realm{}, fmt.Errorf("failed to add vertex: %w", err)
	}

	// TODO
	// emit events to zero
	// - decision
	// - realm creation

	return realm{}, nil
}

func (s *service) createUser(_ context.Context, create userCreate) (user, error) {
	if err := s.g.addVertex(&vertex{
		id:         create.email + "_REALM-USER",
		resource:   "USER",
		edges:      make([]*edge, 0),
		attributes: map[string]interface{}{},
	}); err != nil {
		return user{}, fmt.Errorf("failed to add user vertex: %w", err)
	}

	if err := s.g.addEdge(create.realm, &edge{
		to:           create.email + "_REALM-USER",
		relationship: "RESOURCE",
	}); err != nil {
		return user{}, fmt.Errorf("failed to add user edge to realm vertex: %w", err)
	}

	return user{
		hash: create.email + "_REALM-USER",
	}, nil
}
