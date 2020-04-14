package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
	"github.com/stretchr/testify/assert"
)

func TestCheckGeneratorParams_Validate(t *testing.T) {
	param := &CheckGeneratorParams{}
	assert.NoError(t, validator.Validate(param))

	param = &CheckGeneratorParams{SortBy: "name"}
	assert.NoError(t, validator.Validate(param))

	param = &CheckGeneratorParams{SortBy: "test"}
	assert.Error(t, validator.Validate(param))
}
