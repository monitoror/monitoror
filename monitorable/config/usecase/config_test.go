package usecase

import (
	"errors"
	"testing"

	"github.com/monitoror/monitoror/models/tiles"

	"github.com/monitoror/monitoror/monitorable/ping"
	_pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	"github.com/monitoror/monitoror/monitorable/port"
	_portModels "github.com/monitoror/monitoror/monitorable/port/models"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/monitorable/config/models"

	"github.com/monitoror/monitoror/monitorable/config/mocks"
	. "github.com/stretchr/testify/mock"
)

func initConfigUsecase() *configUsecase {
	usecase := &configUsecase{
		tileConfigs: make(map[tiles.TileType]*TileConfig),
	}

	usecase.RegisterTile(ping.PingTileType, "/ping", &_pingModels.PingParams{})
	usecase.RegisterTile(port.PortTileType, "/port", &_portModels.PortParams{})

	return usecase
}

func TestUsecase_Config_WithUrl_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromUrl", AnythingOfType("string")).Return(&models.Config{}, nil)

	usecase := NewConfigUsecase(mockRepo)

	_, err := usecase.GetConfig(&models.ConfigParams{Url: "test"})
	if assert.NoError(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromUrl", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Config_WithPath_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(&models.Config{}, nil)

	usecase := NewConfigUsecase(mockRepo)

	_, err := usecase.GetConfig(&models.ConfigParams{Path: "test"})
	if assert.NoError(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromPath", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Config_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(nil, errors.New("boom"))

	usecase := NewConfigUsecase(mockRepo)

	_, err := usecase.GetConfig(&models.ConfigParams{Path: "test"})
	if assert.Error(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromPath", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Config_Version(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(&models.Config{}, nil)
	usecase := NewConfigUsecase(mockRepo)

	config, _ := usecase.GetConfig(&models.ConfigParams{Path: "test"})
	assert.Equal(t, CurrentVersion, config.Version)

	mockRepo = new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(&models.Config{Version: 2}, nil)
	usecase = NewConfigUsecase(mockRepo)

	config, _ = usecase.GetConfig(&models.ConfigParams{Path: "test"})
	assert.Equal(t, 2, config.Version)
}
