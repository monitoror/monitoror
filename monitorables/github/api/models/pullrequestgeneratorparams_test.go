package models

import (
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
	"github.com/stretchr/testify/assert"
)

func TestPullRequestGeneratorParams_Validate(t *testing.T) {
	param := &PullRequestGeneratorParams{Owner: "test", Repository: "test"}
	assert.NoError(t, validator.Validate(param))

	param = &PullRequestGeneratorParams{Owner: "test"}
	assert.Error(t, validator.Validate(param))

	param = &PullRequestGeneratorParams{}
	assert.Error(t, validator.Validate(param))
}
