package usecase

import (
	"errors"
	"testing"

	"github.com/monitoror/monitoror/monitorable/jenkins"
	_jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/monitorable/ping"
	_pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	"github.com/monitoror/monitoror/monitorable/port"
	_portModels "github.com/monitoror/monitoror/monitorable/port/models"

	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/config/mocks"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initConfigUsecase(repository config.Repository, conf Cache) *configUsecase {
	usecase := NewConfigUsecase(repository, conf)

	usecase.RegisterTile(ping.PingTileType, &_pingModels.PingParams{}, "/ping")
	usecase.RegisterTile(port.PortTileType, &_portModels.PortParams{}, "/port")
	usecase.RegisterTile(jenkins.JenkinsBuildTileType, &_jenkinsModels.BuildParams{}, "/jenkins/default")

	return usecase.(*configUsecase)
}

func TestUsecase_Load_WithUrl_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromUrl", AnythingOfType("string")).Return(&models.Config{}, nil)

	usecase := initConfigUsecase(mockRepo, Cache{})

	_, err := usecase.GetConfig(&models.ConfigParams{Url: "test"})
	if assert.NoError(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromUrl", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Load_WithPath_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(&models.Config{}, nil)

	usecase := initConfigUsecase(mockRepo, Cache{})

	_, err := usecase.GetConfig(&models.ConfigParams{Path: "test"})
	if assert.NoError(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromPath", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Load_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(nil, errors.New("boom"))

	usecase := initConfigUsecase(mockRepo, Cache{})

	_, err := usecase.GetConfig(&models.ConfigParams{Path: "test"})
	if assert.Error(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromPath", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Load_Version(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(&models.Config{}, nil)
	usecase := initConfigUsecase(mockRepo, Cache{})

	c, _ := usecase.GetConfig(&models.ConfigParams{Path: "test"})
	assert.Equal(t, CurrentVersion, c.Version)

	mockRepo = new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(&models.Config{Version: 2}, nil)
	usecase.repository = mockRepo

	c, _ = usecase.GetConfig(&models.ConfigParams{Path: "test"})
	assert.Equal(t, 2, c.Version)
}
