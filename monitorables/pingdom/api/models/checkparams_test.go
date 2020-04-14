package models

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
	"github.com/stretchr/testify/assert"
)

func TestCheckParams_Validate(t *testing.T) {
	param := &CheckParams{}
	assert.Error(t, validator.Validate(param))

	param = &CheckParams{ID: pointer.ToInt(10)}
	assert.NoError(t, validator.Validate(param))
}
