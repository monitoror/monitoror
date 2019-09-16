package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInfoResponse(t *testing.T) {
	info := NewInfoResponse("a", "b", "c")
	assert.Equal(t, "a", info.Version)
	assert.Equal(t, "b", info.GitCommit)
	assert.Equal(t, "c", info.BuildTime)
}
