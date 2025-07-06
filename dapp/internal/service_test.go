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
			create         = threadCreate{
				name:    threadName,
				members: []string{},
			}
		)

		td, err := svc.createThread(ctx, create)
		assert.Equal(expected, err)
		assert.NotEmpty(td.id)
		assert.Equal(threadName, td.name)
		assert.Equal(create.members, td.members)
		assert.NotEmpty(td.createdAt)
		assert.Nil(td.UpdatedAt)

		t.Logf("successfully created thread %s %s", td.id, td.name)
	})
	t.Run("already_exists", func(t *testing.T) {
		var (
			expected error = errThreadExists
			create         = threadCreate{
				name:    threadName,
				members: []string{},
			}
		)

		_, err := svc.createThread(ctx, create)
		assert.Equal(expected, err)
	})
}
