package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
)

func TestBuildGeneratorParams_Validate(t *testing.T) {
	param := &BuildGeneratorParams{
		Job:     "test",
		Match:   ".*",
		Unmatch: "master",
	}
	assert.NoError(t, validator.Validate(param))

	param = &BuildGeneratorParams{}
	assert.Error(t, validator.Validate(param))

	param = &BuildGeneratorParams{Job: "test", Match: "("}
	assert.Error(t, validator.Validate(param))

	param = &BuildGeneratorParams{Job: "test", Unmatch: "("}
	assert.Error(t, validator.Validate(param))
}
