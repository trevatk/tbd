package identities

import (
	"errors"
	"sync"
)

type edge struct {
	weight       int
	relationship string
	to           string
}

type vertex struct {
	id         string
	resource   string
	edges      []*edge
	attributes map[string]interface{}
}

type graph struct {
	vertices map[string]*vertex
	mu       sync.RWMutex
}

// NewGraph
func NewGraph() *graph {
	return &graph{
		vertices: make(map[string]*vertex, 0),
		mu:       sync.RWMutex{},
	}
}

func bootstrapGraph(g *graph) error {
	permissionVertex := &vertex{
		resource: "PERMISSION",
		edges:    make([]*edge, 0),
		attributes: map[string]interface{}{
			"CREATE_REALM": true,
		},
	}
	err := g.addVertex(permissionVertex)
	if err != nil {
	}

	roleVertex := &vertex{}
	err = g.addVertex(roleVertex)
	if err != nil {
	}

	if err = g.addEdge(roleVertex.id, &edge{
		relationship: "GRANT",
		to:           permissionVertex.id,
	}); err != nil {

	}

	userVertex := &vertex{}
	err = g.addVertex(userVertex)
	if err != nil {
	}

	if err = g.addEdge(userVertex.id, &edge{
		relationship: "ROLE",
		to:           roleVertex.id,
	}); err != nil {
	}

	return nil
}

func (g *graph) addVertex(vertex *vertex) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.vertices[vertex.id]; exists {
		// found existing vertex
		return errors.New("existing vertex found")
	}

	g.vertices[vertex.id] = vertex
	return nil
}

func (g *graph) getVertex(id string) (*vertex, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	v, exists := g.vertices[id]
	if !exists {
		return nil, errors.New("vertex not found")
	}

	return v, nil
}

func (g *graph) addEdge(from string, edge *edge) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	v, exists := g.vertices[from]
	if !exists {
		return errors.New("FROM vertex not found")
	}

	_, exists = g.vertices[edge.to]
	if !exists {
		return errors.New("TO vertex not found")
	}

	v.edges = append(v.edges, edge)
	return nil
}
