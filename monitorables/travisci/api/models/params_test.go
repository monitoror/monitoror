package models

import (
	"fmt"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
	"github.com/stretchr/testify/assert"
)

func TestBuildParams_Validate(t *testing.T) {
	param := &BuildParams{Owner: "test", Repository: "test", Branch: "master"}
	assert.NoError(t, validator.Validate(param))

	param = &BuildParams{Owner: "test", Repository: "test"}
	assert.Error(t, validator.Validate(param))

	param = &BuildParams{Owner: "test"}
	assert.Error(t, validator.Validate(param))

	param = &BuildParams{}
	assert.Error(t, validator.Validate(param))
}

func TestBuildParams_String(t *testing.T) {
	param := &BuildParams{Repository: "test", Owner: "test", Branch: "test"}
	assert.Equal(t, "BUILD-test-test-test", fmt.Sprint(param))
}
