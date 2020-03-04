package usecase

import (
	"errors"
	"testing"

	"github.com/monitoror/monitoror/monitorable/config/mocks"
	"github.com/monitoror/monitoror/monitorable/config/models"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestUsecase_GetConfig_WithURL_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromURL", AnythingOfType("string")).Return(&models.Config{}, nil)

	usecase := initConfigUsecase(mockRepo, nil)

	configBag := usecase.GetConfig(&models.ConfigParams{URL: "test"})
	if assert.Len(t, configBag.Errors, 0) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromURL", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_GetConfig_WithPath_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(&models.Config{}, nil)

	usecase := initConfigUsecase(mockRepo, nil)

	configBag := usecase.GetConfig(&models.ConfigParams{Path: "test"})
	if assert.Len(t, configBag.Errors, 0) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromPath", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_GetConfig_WithError(t *testing.T) {
	for _, testcase := range []struct {
		err     error
		errorID models.ConfigErrorID
	}{
		{
			err:     errors.New("boom"),
			errorID: models.ConfigErrorUnexpectedError,
		},
		{
			err:     &models.ConfigFileNotFoundError{Err: errors.New("boom"), PathOrURL: "path"},
			errorID: models.ConfigErrorConfigNotFound,
		},
		{
			err:     &models.ConfigVersionFormatError{WrongVersion: "18"},
			errorID: models.ConfigErrorUnsupportedVersion,
		},
		{
			err:     &models.ConfigUnmarshalError{Err: errors.New("boom"), RawConfig: ""},
			errorID: models.ConfigErrorUnableToParseConfig,
		},
	} {
		mockRepo := new(mocks.Repository)
		mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(nil, testcase.err)

		usecase := initConfigUsecase(mockRepo, nil)

		configBag := usecase.GetConfig(&models.ConfigParams{Path: "test"})
		if assert.Len(t, configBag.Errors, 1) {
			assert.Equal(t, testcase.errorID, configBag.Errors[0].ID)
			mockRepo.AssertNumberOfCalls(t, "GetConfigFromPath", 1)
			mockRepo.AssertExpectations(t)
		}
	}
}
