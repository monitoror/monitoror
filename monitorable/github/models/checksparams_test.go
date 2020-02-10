package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckParams_IsValid(t *testing.T) {
	p := ChecksParams{Owner: "test", Repository: "test", Ref: "master"}
	assert.True(t, p.IsValid())

	p = ChecksParams{Owner: "test", Repository: "test"}
	assert.False(t, p.IsValid())

	p = ChecksParams{Owner: "test"}
	assert.False(t, p.IsValid())

	p = ChecksParams{}
	assert.False(t, p.IsValid())
}

func TestBuildParams_String(t *testing.T) {
	param := &ChecksParams{Owner: "test", Repository: "test", Ref: "master"}
	assert.Equal(t, "CHECKS-test-test-master", fmt.Sprint(param))
}
