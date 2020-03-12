package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_UnmarshalJSON(t *testing.T) {
	test := &Config{}
	input := `{"version": "1.0"}`
	err := json.Unmarshal([]byte(input), test)
	assert.NoError(t, err)

	input = `{"version": "1.0", "test": "test"}`
	err = json.Unmarshal([]byte(input), test)
	assert.Error(t, err)
	assert.Equal(t, `json: unknown field "test"`, err.Error())
}
