package usecase

import (
	"errors"
	"testing"

	"github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/config/mocks"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	_jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/monitorable/ping"
	_pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	"github.com/monitoror/monitoror/monitorable/port"
	_portModels "github.com/monitoror/monitoror/monitorable/port/models"

	"github.com/AlekSi/pointer"
	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initConfigUsecase(repository config.Repository, store cache.Store) *configUsecase {
	usecase := NewConfigUsecase(repository, store, 5000)

	usecase.RegisterTile(ping.PingTileType, &_pingModels.PingParams{}, "/ping")
	usecase.RegisterTile(port.PortTileType, &_portModels.PortParams{}, "/port")
	usecase.RegisterTile(jenkins.JenkinsBuildTileType, &_jenkinsModels.BuildParams{}, "/jenkins/default")

	return usecase.(*configUsecase)
}

func TestUsecase_Load_WithURL_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromURL", AnythingOfType("string")).Return(&models.Config{}, nil)

	usecase := initConfigUsecase(mockRepo, nil)

	_, err := usecase.GetConfig(&models.ConfigParams{URL: "test"})
	if assert.NoError(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromURL", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Load_WithPath_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(&models.Config{}, nil)

	usecase := initConfigUsecase(mockRepo, nil)

	_, err := usecase.GetConfig(&models.ConfigParams{Path: "test"})
	if assert.NoError(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromPath", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Load_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(nil, errors.New("boom"))

	usecase := initConfigUsecase(mockRepo, nil)

	_, err := usecase.GetConfig(&models.ConfigParams{Path: "test"})
	if assert.Error(t, err) {
		mockRepo.AssertNumberOfCalls(t, "GetConfigFromPath", 1)
		mockRepo.AssertExpectations(t)
	}
}

func TestUsecase_Load_Version(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(&models.Config{}, nil)
	usecase := initConfigUsecase(mockRepo, nil)

	c, _ := usecase.GetConfig(&models.ConfigParams{Path: "test"})
	assert.Nil(t, c.Version)

	mockRepo = new(mocks.Repository)
	mockRepo.On("GetConfigFromPath", AnythingOfType("string")).Return(&models.Config{Version: pointer.ToInt(2)}, nil)
	usecase.repository = mockRepo

	c, _ = usecase.GetConfig(&models.ConfigParams{Path: "test"})
	assert.Equal(t, 2, *c.Version)
}
