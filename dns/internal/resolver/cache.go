package resolver

import (
	"errors"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/trevatk/tbd/lib/protocol/cache/v1"
)

var (
	// ErrKeyNotFound key not found
	ErrKeyNotFound = errors.New("key not found")
)

const (
	defaultInterval = 3
)

// Cache ...
//
//go:generate mockgen -destination mock_cache_test.go -package resolver . Cache
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl int) error

	Start()
	Stop()
}

type memCache struct {
	mu      *sync.RWMutex
	entries map[string]*pb.Entry

	cleanUp  chan struct{}
	interval time.Duration
}

// interface compliance
var _ Cache = (*memCache)(nil)

// NewCache return new in memory cache implementation
func NewCache() Cache {
	return &memCache{
		mu:       &sync.RWMutex{},
		entries:  make(map[string]*pb.Entry),
		cleanUp:  make(chan struct{}),
		interval: defaultInterval,
	}
}

// Start background cache worker
func (m *memCache) Start() {
	go m.worker()
}

// Stop background cache worker
func (m *memCache) Stop() {
	m.cleanUp <- struct{}{}
}

// Get
func (m *memCache) Get(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.entries[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	// verify record has not passed ttl
	// if ttl then expire record in map
	if entry.Ttl != nil {
		expiresAt := entry.Ttl.AsTime()
		if expiresAt.Before(time.Now()) {
			m.mu.Lock()
			defer m.mu.Unlock()

			delete(m.entries, entry.Key)
			return nil, ErrKeyNotFound
		}
	}

	return entry.Value, nil
}

// Set
func (m *memCache) Set(key string, value []byte, ttl int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry := &pb.Entry{
		Key:   key,
		Value: value,
		Ttl:   nil,
	}
	if ttl > numZero {
		ttlAsTime := time.Now().Add(time.Second * time.Duration(ttl))
		entry.Ttl = timestamppb.New(ttlAsTime)
	}

	m.entries[key] = entry

	return nil
}

func (m *memCache) worker() {
	timer := time.NewTicker(m.interval)

	for {
		select {
		case <-m.cleanUp:
			return
		case <-timer.C:
			// iterate over all records
			// compared ttl to current time
			// if ttl is before current time remove from cache
			for key, entry := range m.entries {
				if entry.Ttl != nil {
					expiredAt := entry.Ttl.AsTime()
					if expiredAt.Before(time.Now()) {
						delete(m.entries, key)
					}
				}
			}

			timer.Reset(m.interval)
		}
	}
}
