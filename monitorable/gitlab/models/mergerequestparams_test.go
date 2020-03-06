package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeRequestParams_IsValid(t *testing.T) {
	p := MergeRequestParams{Repository: "test"}
	assert.True(t, p.IsValid())

	p = MergeRequestParams{}
	assert.False(t, p.IsValid())
}
