package nameserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	r1 = &record{
		Domain:     "structx.io",
		RecordType: "A",
		Value:      []byte(host1),
		Ttl:        60,
	}
)

func TestSet(t *testing.T) {
	assert := assert.New(t)

	kv := NewKv()

	t.Run("success", func(t *testing.T) {
		var (
			expected error = nil
		)
		err := kv.set(r1.Domain, r1)
		assert.Equal(expected, err)
	})

	t.Run("key_exists", func(t *testing.T) {
		var (
			expected error = errKeyExists
		)
		err := kv.set(r1.Domain, r1)
		assert.Equal(expected, err)
	})
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	kv := NewKv()
	err := kv.set(r1.Domain, r1)
	assert.NoError(err)

	t.Run("success", func(t *testing.T) {
		var (
			expected error = nil
		)
		value, err := kv.get(r1.Domain)
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
