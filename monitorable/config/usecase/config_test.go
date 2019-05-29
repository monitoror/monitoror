package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/monitorable/config/models"

	"github.com/monitoror/monitoror/monitorable/config/mocks"
	. "github.com/stretchr/testify/mock"
)

func TestUsecase_Config_WithUrl(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromUrl", AnythingOfType("string")).Return(nil, nil)

	usecase := NewConfigUsecase(mockRepo)

	_, err := usecase.Config(&models.ConfigParams{Url: "test"})
	if assert.NoError(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromUrl", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Config_WithPath(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(nil, nil)

	usecase := NewConfigUsecase(mockRepo)

	_, err := usecase.Config(&models.ConfigParams{Path: "test"})
	if assert.NoError(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromPath", 1)
		mockRepo.AssertExpectations(t)
	}
}
