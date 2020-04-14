package models

import (
	"fmt"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"

	"github.com/stretchr/testify/assert"
)

func TestChecksParams_Validate(t *testing.T) {
	param := &ChecksParams{Owner: "test", Repository: "test", Ref: "master"}
	assert.NoError(t, validator.Validate(param))

	param = &ChecksParams{Owner: "test", Repository: "test"}
	assert.Error(t, validator.Validate(param))

	param = &ChecksParams{Owner: "test"}
	assert.Error(t, validator.Validate(param))

	param = &ChecksParams{}
	assert.Error(t, validator.Validate(param))
}

func TestBuildParams_String(t *testing.T) {
	param := &ChecksParams{Owner: "test", Repository: "test", Ref: "master"}
	assert.Equal(t, "CHECKS-test-test-master", fmt.Sprint(param))
}
