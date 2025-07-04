package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalConfig(t *testing.T) {
	cfg := UnmarshalConfig()

	assert.Equal(t, defaultHost, cfg.Gateway.Host)
	assert.Equal(t, defaultPort, cfg.Gateway.Port)

	assert.Equal(t, defaultLogLevel, cfg.Logger.Level)

	assert.Equal(t, defaultSigningKey, cfg.Auth.SigningKey)

	assert.Equal(t, defaultNameserver1, cfg.Nameserver.NS1)
	assert.Equal(t, defaultNameserver2, cfg.Nameserver.NS2)

	assert.Equal(t, defaultKeyValueDir, cfg.KeyValue.Dir)
}
