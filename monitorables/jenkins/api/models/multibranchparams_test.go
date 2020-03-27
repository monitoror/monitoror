package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiBranchParams_IsValid(t *testing.T) {
	p := MultiBranchParams{
		Job:     "test",
		Match:   ".*",
		Unmatch: "master",
	}
	assert.True(t, p.IsValid())

	p = MultiBranchParams{}
	assert.False(t, p.IsValid())

	p = MultiBranchParams{
		Job:   "test",
		Match: "(",
	}
	assert.False(t, p.IsValid())

	p = MultiBranchParams{
		Job:     "test",
		Unmatch: "(",
	}
	assert.False(t, p.IsValid())
}
