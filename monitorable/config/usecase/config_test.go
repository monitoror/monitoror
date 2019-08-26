package usecase

import (
	"errors"
	"testing"

	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config/mocks"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	_jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/monitorable/ping"
	_pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	"github.com/monitoror/monitoror/monitorable/port"
	_portModels "github.com/monitoror/monitoror/monitorable/port/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	mocks2 "github.com/monitoror/monitoror/pkg/monitoror/builder/mocks"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initConfigUsecase() *configUsecase {
	usecase := &configUsecase{
		tileConfigs:        make(map[tiles.TileType]map[string]*TileConfig),
		dynamicTileConfigs: make(map[tiles.TileType]map[string]*DynamicTileConfig),
	}

	params := make(map[string]interface{})
	params["job"] = "test"
	mockBuilder := new(mocks2.DynamicTileBuilder)
	mockBuilder.On("ListDynamicTile", Anything).Return([]builder.Result{{
		TileType: jenkins.JenkinsBuildTileType,
		Params:   params,
	}}, nil)

	mockBuilder2 := new(mocks2.DynamicTileBuilder)
	mockBuilder2.On("ListDynamicTile", Anything).Return(nil, errors.New("boom"))

	usecase.RegisterTile(ping.PingTileType, &_pingModels.PingParams{}, "/ping")
	usecase.RegisterTile(port.PortTileType, &_portModels.PortParams{}, "/port")
	usecase.RegisterTile(jenkins.JenkinsBuildTileType, &_jenkinsModels.BuildParams{}, "/jenkins/default")
	usecase.RegisterTileWithConfigVariant(jenkins.JenkinsBuildTileType, "variant1", &_jenkinsModels.BuildParams{}, "/jenkins/variant1")
	usecase.RegisterDynamicTile(jenkins.JenkinsMultiBranchTileType, &_jenkinsModels.MultiBranchParams{}, mockBuilder)
	usecase.RegisterDynamicTileWithConfigVariant(jenkins.JenkinsMultiBranchTileType, "variant1", &_jenkinsModels.MultiBranchParams{}, mockBuilder2)

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

func TestUsecase_Config_(t *testing.T) {
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
