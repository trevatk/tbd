package internal

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type serviceImpl struct {
	mu      sync.RWMutex
	threads map[string]thread
}

var (
	errThreadExists = errors.New("thread already exists")

	// interface copmliance
	_ service = (*serviceImpl)(nil)
)

// NewService
func NewService() service {
	return &serviceImpl{
		mu:      sync.RWMutex{},
		threads: make(map[string]thread),
	}
}

func (s *serviceImpl) createThread(ctx context.Context, create threadCreate) (thread, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.threads[create.name]; ok {
		return thread{}, errThreadExists
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return thread{}, fmt.Errorf("failed to create uuid: %w", err)
	}

	t := thread{
		id:        uid.String(),
		name:      create.name,
		members:   create.members,
		createdAt: time.Now(),
		UpdatedAt: nil,
	}

	s.threads[create.name] = t

	return t, nil
}
