package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultError(t *testing.T) {
	err := NewDefaultError("Test", "test > 23")
	assert.Equal(t, ErrorDefault, err.GetErrorID())
	assert.Equal(t, "test > 23", err.Expected())
	assert.Equal(t, `Invalid "Test" field. Must be test > 23.`, err.Error())

	assert.Equal(t, "Test", err.GetFieldName())
	err.SetFieldName("Test2")
	assert.Equal(t, "Test2", err.GetFieldName())

	err = NewDefaultError("Test", "")
	assert.Equal(t, "", err.Expected())
	assert.Equal(t, `Invalid "Test" field.`, err.Error())
}
