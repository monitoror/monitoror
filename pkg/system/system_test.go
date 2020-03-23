package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRawSocketAvailable(t *testing.T) {
	// Can't test this one better.
	assert.NotPanics(t, func() { IsRawSocketAvailable() })
}

func TestListLocalhostIpv4(t *testing.T) {
	ips, err := ListLocalhostIpv4()
	assert.NoError(t, err)
	assert.Contains(t, ips, "127.0.0.1")
}
