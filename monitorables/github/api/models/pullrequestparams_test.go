package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPullRequestParams_IsValid(t *testing.T) {
	p := PullRequestParams{Owner: "test", Repository: "test"}
	assert.True(t, p.IsValid())

	p = PullRequestParams{Owner: "test"}
	assert.False(t, p.IsValid())

	p = PullRequestParams{}
	assert.False(t, p.IsValid())
}
