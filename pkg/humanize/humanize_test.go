package humanize

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestGetHumanizeInterface(t *testing.T) {
	assert.Equal(t, "100000000", Interface(float64(100000000)))
	assert.Equal(t, "100000000", Interface(100000000))
	assert.Equal(t, "100000000", Interface("100000000"))
	assert.Equal(t, "aaa", Interface("aaa"))
	assert.Equal(t, "aaa", Interface(pointer.ToString("aaa")))

	var test *string
	assert.Equal(t, "", Interface(test))

	// TODO handle slice properly if needed
	assert.Equal(t, "[1e+08]", Interface([]float64{100000000}))
}
