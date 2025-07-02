package authoritative

import (
	"errors"
	"sync"
)

type record struct {
	domain     string
	recordType string // CNAME, A, MX etc...
	value      []byte // IP address, CNAME value
	ttl        int64
}

//go:generate mockgen -destination mock_kv_test.go -package authoritative . kv
type kv interface {
	get(string) (*record, error)
	set(string, *record) error
}

type inMemoryKv struct {
	mu     sync.RWMutex
	values map[string]*record
}

// interface compliance
var _ kv = (*inMemoryKv)(nil)

// NewKv return new key value store implementation
func NewKv() kv {
	return &inMemoryKv{
		mu:     sync.RWMutex{},
		values: make(map[string]*record),
	}
}

func (k *inMemoryKv) get(key string) (*record, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	value, ok := k.values[key]
	if !ok {
		return nil, errKeyNotFound
	}

	return value, nil
}

func (k *inMemoryKv) set(key string, value *record) error {
	if value == nil {
		return errors.New("nil value")
	}

	k.mu.Lock()
	defer k.mu.Unlock()

	_, ok := k.values[key]
	if ok {
		return errKeyExists
	}

	k.values[key] = value
	return nil
}
