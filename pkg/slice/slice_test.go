package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	i, find := Find([]string{"TEST", "TEST2"}, "TEST")
	assert.Equal(t, 0, i)
	assert.True(t, find)

	_, find = Find([]string{"TEST", "TEST2"}, "TEST3")
	assert.False(t, find)
}
