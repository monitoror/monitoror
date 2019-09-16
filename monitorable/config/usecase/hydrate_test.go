package usecase

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/config/repository"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	_jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	. "github.com/monitoror/monitoror/pkg/monitoror/builder/mocks"
	mocks2 "github.com/monitoror/monitoror/pkg/monitoror/builder/mocks"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func TestUsecase_Hydrate(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "EMPTY" },
    { "type": "PING", "params": { "hostname": "aserver.com" } },
    { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } },
    { "type": "GROUP", "label": "...", "tiles": [
      { "type": "PING", "params": { "hostname": "aserver.com" } },
      { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } }
    ]},
		{ "type": "JENKINS-BUILD", "params": { "job": "test" } },
		{ "type": "JENKINS-BUILD", "configVariant": "variant1", "params": { "job": "test" } }
  ]
}
`

	usecase := initConfigUsecase(nil, Cache{Expire: 1000, CleanupInterval: 1000})
	usecase.RegisterTileWithConfigVariant(jenkins.JenkinsBuildTileType, "variant1", &_jenkinsModels.BuildParams{}, "/jenkins/variant1")

	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.ReadConfig(reader)
	assert.NoError(t, err)

	usecase.Hydrate(config, "http://localhost:8080")
	assert.Len(t, config.Errors, 0)
	assert.Len(t, config.Warnings, 0)

	assert.Equal(t, "http://localhost:8080/ping?hostname=aserver.com", config.Tiles[1].Url)
	assert.Equal(t, "http://localhost:8080/port?hostname=bserver.com&port=22", config.Tiles[2].Url)

	group := config.Tiles[3].Tiles
	assert.Equal(t, "http://localhost:8080/ping?hostname=aserver.com", group[0].Url)
	assert.Equal(t, "http://localhost:8080/port?hostname=bserver.com&port=22", group[1].Url)

	assert.Equal(t, "http://localhost:8080/jenkins/default?job=test", config.Tiles[4].Url)
	assert.Equal(t, "http://localhost:8080/jenkins/variant1?job=test", config.Tiles[5].Url)
}

func TestUsecase_Hydrate_WithDynamic(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "JENKINS-MULTIBRANCH", "params": {"job": "test"}},
    { "type": "GROUP", "label": "...", "tiles": [
      { "type": "PING", "params": { "hostname": "aserver.com" } },
			{ "type": "JENKINS-MULTIBRANCH", "params": {"job": "test"}}
    ]},
    { "type": "GROUP", "label": "...", "tiles": [
    	{ "type": "JENKINS-MULTIBRANCH", "params": {"job": "test"}}
    ]}
  ]
}
`
	params := make(map[string]interface{})
	params["job"] = "test"
	mockBuilder := new(DynamicTileBuilder)
	mockBuilder.On("ListDynamicTile", Anything).Return([]builder.Result{{TileType: jenkins.JenkinsBuildTileType, Params: params}}, nil)

	usecase := initConfigUsecase(nil, Cache{Expire: 1000, CleanupInterval: 1000})
	usecase.RegisterDynamicTile(jenkins.JenkinsMultiBranchTileType, &_jenkinsModels.MultiBranchParams{}, mockBuilder)

	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.ReadConfig(reader)
	assert.NoError(t, err)

	usecase.Hydrate(config, "http://localhost:8080")
	assert.Len(t, config.Errors, 0)
	assert.Len(t, config.Warnings, 0)

	assert.Equal(t, 3, len(config.Tiles))
	assert.Equal(t, jenkins.JenkinsBuildTileType, config.Tiles[0].Type)
	assert.Equal(t, "http://localhost:8080/jenkins/default?job=test", config.Tiles[0].Url)
	assert.Equal(t, jenkins.JenkinsBuildTileType, config.Tiles[1].Tiles[1].Type)
	assert.Equal(t, "http://localhost:8080/jenkins/default?job=test", config.Tiles[1].Tiles[1].Url)
	assert.Equal(t, jenkins.JenkinsBuildTileType, config.Tiles[2].Tiles[0].Type)
	assert.Equal(t, "http://localhost:8080/jenkins/default?job=test", config.Tiles[2].Tiles[0].Url)
	mockBuilder.AssertNumberOfCalls(t, "ListDynamicTile", 3)
	mockBuilder.AssertExpectations(t)
}

func TestUsecase_Hydrate_WithDynamic_WithError(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "JENKINS-MULTIBRANCH", "params": {"job": "test"}},
    { "type": "GROUP", "label": "...", "tiles": [
      { "type": "PING", "params": { "hostname": "aserver.com" } },
			{ "type": "JENKINS-MULTIBRANCH", "params": {"job": "test"}}
    ]},
    { "type": "GROUP", "label": "...", "tiles": [
    	{ "type": "JENKINS-MULTIBRANCH", "configVariant": "variant1", "params": {"job": "test"}}
    ]}
  ]
}
`
	params := make(map[string]interface{})
	params["job"] = "test"
	mockBuilder := new(DynamicTileBuilder)
	mockBuilder.On("ListDynamicTile", Anything).Return([]builder.Result{{TileType: jenkins.JenkinsBuildTileType, Params: params}}, nil)
	mockBuilder2 := new(mocks2.DynamicTileBuilder)
	mockBuilder2.On("ListDynamicTile", Anything).Return(nil, errors.New("unable to found job"))

	usecase := initConfigUsecase(nil, Cache{Expire: 1000, CleanupInterval: 1000})
	usecase.RegisterTileWithConfigVariant(jenkins.JenkinsBuildTileType, "variant1", &_jenkinsModels.BuildParams{}, "/jenkins/variant1")
	usecase.RegisterDynamicTile(jenkins.JenkinsMultiBranchTileType, &_jenkinsModels.MultiBranchParams{}, mockBuilder)
	usecase.RegisterDynamicTileWithConfigVariant(jenkins.JenkinsMultiBranchTileType, "variant1", &_jenkinsModels.MultiBranchParams{}, mockBuilder2)

	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.ReadConfig(reader)
	assert.NoError(t, err)

	usecase.Hydrate(config, "http://localhost:8080")
	assert.Len(t, config.Errors, 1)
	assert.Len(t, config.Warnings, 0)
	assert.Contains(t, config.Errors, `Error while listing JENKINS-MULTIBRANCH dynamic tiles (params: {"job":"test"}). unable to found job`)
	mockBuilder.AssertNumberOfCalls(t, "ListDynamicTile", 2)
	mockBuilder.AssertExpectations(t)
	mockBuilder2.AssertNumberOfCalls(t, "ListDynamicTile", 1)
	mockBuilder2.AssertExpectations(t)
}

func TestUsecase_Hydrate_WithDynamic_WithWarning(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "JENKINS-MULTIBRANCH", "params": {"job": "test"}},
    { "type": "GROUP", "label": "...", "tiles": [
      { "type": "PING", "params": { "hostname": "aserver.com" } },
			{ "type": "JENKINS-MULTIBRANCH", "params": {"job": "test"}}
    ]},
    { "type": "GROUP", "label": "...", "tiles": [
    	{ "type": "JENKINS-MULTIBRANCH", "configVariant": "variant1", "params": {"job": "test"}}
    ]}
  ]
}
`
	params := make(map[string]interface{})
	params["job"] = "test"
	mockBuilder := new(DynamicTileBuilder)
	mockBuilder.On("ListDynamicTile", Anything).Return([]builder.Result{{TileType: jenkins.JenkinsBuildTileType, Params: params}}, nil)
	mockBuilder2 := new(mocks2.DynamicTileBuilder)
	mockBuilder2.On("ListDynamicTile", Anything).Return(nil, context.DeadlineExceeded)

	usecase := initConfigUsecase(nil, Cache{Expire: 1000, CleanupInterval: 1000})
	usecase.RegisterTileWithConfigVariant(jenkins.JenkinsBuildTileType, "variant1", &_jenkinsModels.BuildParams{}, "/jenkins/variant1")
	usecase.RegisterDynamicTile(jenkins.JenkinsMultiBranchTileType, &_jenkinsModels.MultiBranchParams{}, mockBuilder)
	usecase.RegisterDynamicTileWithConfigVariant(jenkins.JenkinsMultiBranchTileType, "variant1", &_jenkinsModels.MultiBranchParams{}, mockBuilder2)

	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.ReadConfig(reader)
	assert.NoError(t, err)

	usecase.Hydrate(config, "http://localhost:8080")
	assert.Len(t, config.Errors, 0)
	assert.Len(t, config.Warnings, 1)
	assert.Contains(t, config.Warnings, `Warning while listing JENKINS-MULTIBRANCH dynamic tiles (params: {"job":"test"}). timeout/host unreachable`)
	mockBuilder.AssertNumberOfCalls(t, "ListDynamicTile", 2)
	mockBuilder.AssertExpectations(t)
	mockBuilder2.AssertNumberOfCalls(t, "ListDynamicTile", 1)
	mockBuilder2.AssertExpectations(t)
}

func TestUsecase_Hydrate_WithDynamic_WithTimeoutCache(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "JENKINS-MULTIBRANCH", "params": {"job": "test"}}
	]
}
`
	usecase := initConfigUsecase(nil, Cache{Expire: 1000, CleanupInterval: 1000})

	params := make(map[string]interface{})
	params["job"] = "test"
	cachedResult := []builder.Result{{TileType: jenkins.JenkinsBuildTileType, Params: params}}
	cacheKey := fmt.Sprintf("%s_%s_%s", "JENKINS-MULTIBRANCH", "default", `{"job":"test"}`)
	_ = usecase.dynamicTileCache.Add(cacheKey, cachedResult, 0)

	mockBuilder := new(DynamicTileBuilder)
	mockBuilder.On("ListDynamicTile", Anything).Return(nil, context.DeadlineExceeded)
	usecase.RegisterDynamicTile(jenkins.JenkinsMultiBranchTileType, &_jenkinsModels.MultiBranchParams{}, mockBuilder)

	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.ReadConfig(reader)
	if assert.NoError(t, err) {
		usecase.Hydrate(config, "http://localhost:8080")
		assert.Len(t, config.Errors, 0)
		assert.Len(t, config.Warnings, 0)
		assert.Equal(t, jenkins.JenkinsBuildTileType, config.Tiles[0].Type)
		assert.Equal(t, "http://localhost:8080/jenkins/default?job=test", config.Tiles[0].Url)
		mockBuilder.AssertNumberOfCalls(t, "ListDynamicTile", 1)
		mockBuilder.AssertExpectations(t)
	}
}
