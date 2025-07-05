package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	threadName = "helloworld"
)

func TestCreateThread(t *testing.T) {

	ctx := t.Context()

	assert := assert.New(t)

	svc := NewService()
	t.Run("new_thread", func(t *testing.T) {
		var (
			expected error = nil
		)

		td, err := svc.createThread(ctx, threadName)
		assert.Equal(expected, err)
		assert.NotEmpty(td.id)
		assert.Equal(threadName, td.name)

		t.Logf("successfully created thread %s %s", td.id, td.name)
	})
	t.Run("already_exists", func(t *testing.T) {
		var (
			expected error = errThreadExists
		)

		_, err := svc.createThread(ctx, threadName)
		assert.Equal(expected, err)
	})
}
