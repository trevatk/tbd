package authoritative

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	r1 = &record{
		domain:     "structx.io",
		recordType: "A",
		value:      []byte(host1),
		ttl:        60,
	}
)

func TestSet(t *testing.T) {
	assert := assert.New(t)

	kv := NewKv()

	t.Run("success", func(t *testing.T) {
		var (
			expected error = nil
		)
		err := kv.set(r1.domain, r1)
		assert.Equal(expected, err)
	})

	t.Run("key_exists", func(t *testing.T) {
		var (
			expected error = errKeyExists
		)
		err := kv.set(r1.domain, r1)
		assert.Equal(expected, err)
	})
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	kv := NewKv()
	err := kv.set(r1.domain, r1)
	assert.NoError(err)

	t.Run("success", func(t *testing.T) {
		var (
			expected error = nil
		)
		value, err := kv.get(r1.domain)
		assert.Equal(expected, err)
		assert.Equal(r1, value)
	})

	t.Run("key_not_exists", func(t *testing.T) {
		var (
			expected error = errKeyNotFound
		)
		_, err := kv.get("not_found")
		assert.Equal(expected, err)
	})
}
