package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIssuesParams_IsValid(t *testing.T) {
	p := IssuesParams{Query: "test"}
	assert.True(t, p.IsValid())

	p = IssuesParams{}
	assert.False(t, p.IsValid())
}
