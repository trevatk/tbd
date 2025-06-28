package authoritative

import (
	"errors"
	"sync"
)

var (
	errKeyNotFound = errors.New("key not found")
)

type record struct {
	domain     string
	recordType string // CNAME, A, MX etc...
	value      string // IP address, CNAME value
	ttl        int64
}

// Kv
//
//go:generate mockgen -destination mock_kv_test.go -package authoritative . Kv
type Kv interface {
	get(string) (record, error)
	set(string, record) error
}

type inMemoryKv struct {
	mu     *sync.RWMutex
	values map[string]record
}

// interface compliance
var _ Kv = (*inMemoryKv)(nil)

func NewKv() Kv {
	return &inMemoryKv{}
}

// get
func (k *inMemoryKv) get(key string) (record, error) {
	k.mu.RLock()
	defer k.mu.Unlock()

	value, ok := k.values[key]
	if !ok {
		return record{}, errKeyNotFound
	}

	return value, nil
}

// set
func (k *inMemoryKv) set(key string, value record) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.values[key] = value
	return nil
}
