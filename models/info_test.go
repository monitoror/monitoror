package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInfoResponse(t *testing.T) {
	info := NewInfoResponse("version", "gitCommit", "buildTime", "buildTags")
	assert.Equal(t, "version", info.Version)
	assert.Equal(t, "gitCommit", info.GitCommit)
	assert.Equal(t, "buildTime", info.BuildTime)
	assert.Equal(t, "buildTags", info.BuildTags)
}
