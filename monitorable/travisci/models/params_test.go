package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildParams_IsValid(t *testing.T) {
	param := &BuildParams{}
	assert.False(t, param.IsValid())

	param = &BuildParams{Repository: "test"}
	assert.False(t, param.IsValid())

	param = &BuildParams{Repository: "test", Owner: "test"}
	assert.False(t, param.IsValid())

	param = &BuildParams{Repository: "test", Owner: "test", Branch: "test"}
	assert.True(t, param.IsValid())
}

func TestBuildParams_String(t *testing.T) {
	param := &BuildParams{Repository: "test", Owner: "test", Branch: "test"}
	assert.Equal(t, "BUILD-test-test-test", fmt.Sprint(param))
}
