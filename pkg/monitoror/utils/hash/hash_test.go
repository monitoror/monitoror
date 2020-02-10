package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMD5Hash(t *testing.T) {
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", GetMD5Hash("test"))
}
