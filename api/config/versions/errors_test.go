package versions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigVersionFormatError(t *testing.T) {
	err := &ConfigVersionFormatError{WrongVersion: "10"}
	assert.Equal(t, "json: cannot unmarshal 10 into Go struct field Config.Version of type string and X.y format", err.Error())
}
