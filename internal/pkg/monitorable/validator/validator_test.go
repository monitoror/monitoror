package validator

import (
	"testing"

	"github.com/monitoror/monitoror/api/config/mocks"
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidate(t *testing.T) {
	mockValidator := new(mocks.ParamsValidator)
	mockValidator.On("Validate", mock.Anything).Return(nil)
	assert.NoError(t, Validate(mockValidator))

	mockValidator2 := new(mocks.ParamsValidator)
	mockValidator2.On("Validate", mock.Anything).Return(&uiConfigModels.ConfigError{Message: "boom"})
	err := Validate(mockValidator2)
	assert.Error(t, err)
	assert.Equal(t, "boom", err.Error())
}
