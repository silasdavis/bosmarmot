package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitKeyClient(t *testing.T) {
	localKeyClient, err := InitKeyClient(DefaultKeysURL())
	require.NoError(t, err)
	err = localKeyClient.HealthCheck()
	assert.NoError(t, err)
}
