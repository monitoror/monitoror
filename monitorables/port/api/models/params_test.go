package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
	"github.com/stretchr/testify/assert"
)

func TestPortParams_Validate(t *testing.T) {
	param := &PortParams{}
	assert.Error(t, validator.Validate(param))

	param = &PortParams{Hostname: "test"}
	assert.Error(t, validator.Validate(param))

	param = &PortParams{Hostname: "test", Port: 22}
	assert.NoError(t, validator.Validate(param))
}
