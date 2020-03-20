package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountParams_IsValid(t *testing.T) {
	p := CountParams{Query: "test"}
	assert.True(t, p.IsValid())

	p = CountParams{}
	assert.False(t, p.IsValid())
}
