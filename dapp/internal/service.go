package internal

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var (
	errThreadExists = errors.New("thread already exists")
)

type serviceImpl struct {
	mu      sync.RWMutex
	threads map[string]thread
}

// interface compliance
var _ service = (*serviceImpl)(nil)

// NewService
func NewService() service {
	return &serviceImpl{
		mu:      sync.RWMutex{},
		threads: make(map[string]thread),
	}
}

func (s *serviceImpl) createThread(ctx context.Context, name string) (thread, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.threads[name]; ok {
		return thread{}, errThreadExists
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return thread{}, fmt.Errorf("failed to create uuid: %w", err)
	}

	t := thread{
		id:   uid.String(),
		name: name,
	}

	s.threads[name] = t

	return t, nil
}
