package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildParams_IsValid(t *testing.T) {
	p := BuildParams{
		Job:    "test",
		Branch: "test",
	}
	assert.True(t, p.IsValid())

	p = BuildParams{}
	assert.False(t, p.IsValid())
}
