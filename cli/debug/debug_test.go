package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	Disable()
	assert.False(t, IsEnabled())
	Enable()
	assert.True(t, IsEnabled())
}
