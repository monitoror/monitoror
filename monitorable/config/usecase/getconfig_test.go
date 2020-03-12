package usecase

import (
	"errors"
	"fmt"
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
		err       error
		errorID   models.ConfigErrorID
		errorData models.ConfigErrorData
	}{
		{
			err:     errors.New("boom"),
			errorID: models.ConfigErrorUnexpectedError,
		},
		{
			err:       &models.ConfigFileNotFoundError{Err: errors.New("boom"), PathOrURL: "path"},
			errorID:   models.ConfigErrorConfigNotFound,
			errorData: models.ConfigErrorData{Value: "path"},
		},
		{
			err:     &models.ConfigVersionFormatError{WrongVersion: "18"},
			errorID: models.ConfigErrorUnsupportedVersion,
			errorData: models.ConfigErrorData{
				Value:     "18",
				FieldName: "version",
				Expected:  fmt.Sprintf("%q >= version >= %q", MinimalVersion, CurrentVersion),
			},
		},
		{
			err:       &models.ConfigUnmarshalError{Err: errors.New("boom"), RawConfig: "test json"},
			errorID:   models.ConfigErrorUnableToParseConfig,
			errorData: models.ConfigErrorData{ConfigExtract: "test json"},
		},
		{
			err:       &models.ConfigUnmarshalError{Err: errors.New(`json: unknown field "test"`), RawConfig: "test json"},
			errorID:   models.ConfigErrorUnknownField,
			errorData: models.ConfigErrorData{FieldName: "test", ConfigExtract: "test json"},
		},
	} {
		mockRepo := new(mocks.Repository)
		mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(nil, testcase.err)

		usecase := initConfigUsecase(mockRepo, nil)

		configBag := usecase.GetConfig(&models.ConfigParams{Path: "test"})
		if assert.Len(t, configBag.Errors, 1) {
			assert.Equal(t, testcase.errorID, configBag.Errors[0].ID)
			assert.Equal(t, testcase.errorData, configBag.Errors[0].Data)
			mockRepo.AssertNumberOfCalls(t, "GetConfigFromPath", 1)
			mockRepo.AssertExpectations(t)
		}
	}
}
