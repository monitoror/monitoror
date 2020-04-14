package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
	"github.com/stretchr/testify/assert"
)

func TestPingParams_Validate(t *testing.T) {
	param := &PingParams{Hostname: "test"}
	assert.NoError(t, validator.Validate(param))

	param = &PingParams{}
	assert.Error(t, validator.Validate(param))
}
