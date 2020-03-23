package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRawSocketAvailable(t *testing.T) {
	// Can't test this one better.
	assert.NotPanics(t, func() { IsRawSocketAvailable() })
}
