package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trevatk/tbd/lib/logging"
	v1 "github.com/trevatk/tbd/lib/protocol/wellknown/v1"
	"github.com/trevatk/tbd/lib/setup"
)

func TestGetDIDConfiguration(t *testing.T) {

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	assert := assert.New(t)

	cfg := setup.UnmarshalConfig()
	logger := logging.New(cfg.Logger.Level)
	g := newGrpcTransport(logger, "testfiles/structx.local.wellknown")
	t.Run("success", func(t *testing.T) {
		var (
			expected error = nil
		)
		_, err := g.GetDIDConfiguration(ctx, &v1.GetDIDConfigurationRequest{})
		assert.Equal(expected, err)
	})
}
