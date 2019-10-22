package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildParams_IsValid(t *testing.T) {
	p := BuildParams{Job: "test", Branch: "test"}
	assert.True(t, p.IsValid())

	p = BuildParams{Job: "test"}
	assert.True(t, p.IsValid())

	p = BuildParams{}
	assert.False(t, p.IsValid())
}

func TestBuildParams_String(t *testing.T) {
	param := &BuildParams{Job: "test", Branch: "test"}
	assert.Equal(t, "BUILD-test-test", fmt.Sprint(param))
}
