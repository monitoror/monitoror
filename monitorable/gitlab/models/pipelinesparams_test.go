package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipelinesParams_IsValid(t *testing.T) {
	p := PipelinesParams{Repository: "test", Ref: "master"}
	assert.True(t, p.IsValid())

	p = PipelinesParams{Repository: "test"}
	assert.False(t, p.IsValid())

	p = PipelinesParams{Ref: "master"}
	assert.False(t, p.IsValid())

	p = PipelinesParams{}
	assert.False(t, p.IsValid())
}

func TestPipelinesParams_String(t *testing.T) {
	param := &PipelinesParams{Repository: "test", Ref: "master"}
	assert.Equal(t, "PIPELINES-test-master", fmt.Sprint(param))
}
