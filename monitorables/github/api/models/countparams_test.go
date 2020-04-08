package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
	"github.com/stretchr/testify/assert"
)

func TestCountParams_Validate(t *testing.T) {
	param := &CountParams{Query: "test"}
	assert.NoError(t, validator.Validate(param))

	param = &CountParams{}
	assert.Error(t, validator.Validate(param))
}
