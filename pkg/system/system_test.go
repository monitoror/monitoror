package system

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRawSocketAvailable(t *testing.T) {
	// Can't test this one better.
	assert.NotPanics(t, func() { IsRawSocketAvailable() })
}

func TestGetNetworkIp(t *testing.T) {
	ip := GetNetworkIP()
	fmt.Println(ip)
	assert.NotEmpty(t, ip)
}
