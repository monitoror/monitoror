package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	jenkinsApi "github.com/monitoror/monitoror/monitorables/jenkins/api"
	jenkinsModels "github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/stretchr/testify/assert"
)

func TestUsecase_Hydrate(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "EMPTY" },
    { "type": "PING", "params": { "hostname": "aserver.com", "values": [123, 456] } },
    { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } },
    { "type": "GROUP", "label": "...", "tiles": [
      { "type": "PING", "params": { "hostname": "aserver.com" } },
      { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } }
    ]},
		{ "type": "JENKINS-BUILD", "params": { "job": "test" } },
		{ "type": "JENKINS-BUILD", "configVariant": "variant1", "params": { "job": "test" } },
    { "type": "PINGDOM-CHECK", "params": { "id": 10000000 } }
  ]
}
`

	store := cache.NewGoCacheStore(time.Second, time.Second)
	usecase := initConfigUsecase(nil, store)
	usecase.EnableTile(jenkinsApi.JenkinsBuildTileType, "variant1", &jenkinsModels.BuildParams{}, "/jenkins/variant1", 1000)

	config, err := readConfig(input)

	assert.NoError(t, err)

	usecase.Hydrate(config)
	assert.Len(t, config.Errors, 0)

	assert.Equal(t, "/ping?hostname=aserver.com&values=123&values=456", config.Config.Tiles[1].URL)
	assert.Equal(t, 1000, *config.Config.Tiles[1].InitialMaxDelay)
	assert.Equal(t, "/port?hostname=bserver.com&port=22", config.Config.Tiles[2].URL)
	assert.Equal(t, 1000, *config.Config.Tiles[2].InitialMaxDelay)

	group := config.Config.Tiles[3].Tiles
	assert.Equal(t, "/ping?hostname=aserver.com", group[0].URL)
	assert.Equal(t, 1000, *group[0].InitialMaxDelay)
	assert.Equal(t, "/port?hostname=bserver.com&port=22", group[1].URL)
	assert.Equal(t, 1000, *group[1].InitialMaxDelay)

	assert.Equal(t, "/jenkins/default?job=test", config.Config.Tiles[4].URL)
	assert.Equal(t, 1000, *config.Config.Tiles[4].InitialMaxDelay)
	assert.Equal(t, "/jenkins/variant1?job=test", config.Config.Tiles[5].URL)
	assert.Equal(t, 1000, *config.Config.Tiles[5].InitialMaxDelay)
	assert.Equal(t, "/pingdom/default?id=10000000", config.Config.Tiles[6].URL)
	assert.Equal(t, 1000, *config.Config.Tiles[6].InitialMaxDelay)
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
	mockBuilder := func(_ interface{}) ([]builder.Result, error) {
		return []builder.Result{{TileType: jenkinsApi.JenkinsBuildTileType, Params: params}}, nil
	}

	store := cache.NewGoCacheStore(time.Second, time.Second)
	usecase := initConfigUsecase(nil, store)
	usecase.EnableDynamicTile(jenkinsApi.JenkinsMultiBranchTileType, coreModels.DefaultVariant, &jenkinsModels.MultiBranchParams{}, mockBuilder)

	config, err := readConfig(input)
	assert.NoError(t, err)

	usecase.Hydrate(config)
	assert.Len(t, config.Errors, 0)

	assert.Equal(t, 3, len(config.Config.Tiles))
	assert.Equal(t, jenkinsApi.JenkinsBuildTileType, config.Config.Tiles[0].Type)
	assert.Equal(t, "/jenkins/default?job=test", config.Config.Tiles[0].URL)
	assert.Equal(t, 1000, *config.Config.Tiles[0].InitialMaxDelay)
	assert.Equal(t, jenkinsApi.JenkinsBuildTileType, config.Config.Tiles[1].Tiles[1].Type)
	assert.Equal(t, "/jenkins/default?job=test", config.Config.Tiles[1].Tiles[1].URL)
	assert.Equal(t, 1000, *config.Config.Tiles[1].Tiles[1].InitialMaxDelay)
	assert.Equal(t, jenkinsApi.JenkinsBuildTileType, config.Config.Tiles[2].Tiles[0].Type)
	assert.Equal(t, "/jenkins/default?job=test", config.Config.Tiles[2].Tiles[0].URL)
	assert.Equal(t, 1000, *config.Config.Tiles[2].Tiles[0].InitialMaxDelay)
}

func TestUsecase_Hydrate_WithDynamicEmpty(t *testing.T) {
	input := `
{
  "columns": 4,
  "tiles": [
    { "type": "PING", "params": { "hostname": "aserver.com" } },
    { "type": "GROUP", "label": "...", "tiles": [
    	{ "type": "JENKINS-MULTIBRANCH", "params": {"job": "test"}}
    ]},
    { "type": "PING", "params": { "hostname": "bserver.com" } }
  ]
}
`
	params := make(map[string]interface{})
	params["job"] = "test"
	mockBuilder := func(_ interface{}) ([]builder.Result, error) { return []builder.Result{}, nil }

	store := cache.NewGoCacheStore(time.Second, time.Second)
	usecase := initConfigUsecase(nil, store)
	usecase.EnableDynamicTile(jenkinsApi.JenkinsMultiBranchTileType, coreModels.DefaultVariant, &jenkinsModels.MultiBranchParams{}, mockBuilder)

	config, err := readConfig(input)
	assert.NoError(t, err)

	usecase.Hydrate(config)
	assert.Len(t, config.Errors, 0)

	assert.Equal(t, 2, len(config.Config.Tiles))
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
	mockBuilder := func(_ interface{}) ([]builder.Result, error) {
		return []builder.Result{{TileType: jenkinsApi.JenkinsBuildTileType, Params: params}}, nil
	}
	mockBuilder2 := func(_ interface{}) ([]builder.Result, error) { return nil, errors.New("unable to find job") }

	store := cache.NewGoCacheStore(time.Second, time.Second)
	usecase := initConfigUsecase(nil, store)
	usecase.EnableTile(jenkinsApi.JenkinsBuildTileType, "variant1", &jenkinsModels.BuildParams{}, "/jenkins/variant1", 1000)
	usecase.EnableDynamicTile(jenkinsApi.JenkinsMultiBranchTileType, coreModels.DefaultVariant, &jenkinsModels.MultiBranchParams{}, mockBuilder)
	usecase.EnableDynamicTile(jenkinsApi.JenkinsMultiBranchTileType, "variant1", &jenkinsModels.MultiBranchParams{}, mockBuilder2)

	config, err := readConfig(input)
	assert.NoError(t, err)

	usecase.Hydrate(config)
	assert.Len(t, config.Errors, 1)
	assert.Equal(t, config.Errors[0].ID, models.ConfigErrorUnableToHydrate)
	assert.Contains(t, config.Errors[0].Data.ConfigExtract, `JENKINS-MULTIBRANCH`)
}

func TestUsecase_Hydrate_WithDynamic_WithTimeoutError(t *testing.T) {
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
	mockBuilder := func(_ interface{}) ([]builder.Result, error) {
		return []builder.Result{{TileType: jenkinsApi.JenkinsBuildTileType, Params: params}}, nil
	}
	mockBuilder2 := func(_ interface{}) ([]builder.Result, error) { return nil, context.DeadlineExceeded }

	store := cache.NewGoCacheStore(time.Second, time.Second)
	usecase := initConfigUsecase(nil, store)
	usecase.EnableTile(jenkinsApi.JenkinsBuildTileType, "variant1", &jenkinsModels.BuildParams{}, "/jenkins/variant1", 1000)
	usecase.EnableDynamicTile(jenkinsApi.JenkinsMultiBranchTileType, coreModels.DefaultVariant, &jenkinsModels.MultiBranchParams{}, mockBuilder)
	usecase.EnableDynamicTile(jenkinsApi.JenkinsMultiBranchTileType, "variant1", &jenkinsModels.MultiBranchParams{}, mockBuilder2)

	config, err := readConfig(input)
	assert.NoError(t, err)

	usecase.Hydrate(config)
	assert.Len(t, config.Errors, 1)
	assert.Equal(t, config.Errors[0].ID, models.ConfigErrorUnableToHydrate)
	assert.Contains(t, config.Errors[0].Data.ConfigExtract, `JENKINS-MULTIBRANCH`)
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
	store := cache.NewGoCacheStore(time.Second, time.Second)
	usecase := initConfigUsecase(nil, store)

	params := make(map[string]interface{})
	params["job"] = "test"
	cachedResult := []builder.Result{{TileType: jenkinsApi.JenkinsBuildTileType, Params: params}}
	cacheKey := fmt.Sprintf("%s:%s_%s_%s", DynamicTileStoreKeyPrefix, "JENKINS-MULTIBRANCH", "default", `{"job":"test"}`)
	_ = usecase.dynamicTileStore.Add(cacheKey, cachedResult, 0)

	mockBuilder := func(_ interface{}) ([]builder.Result, error) { return nil, context.DeadlineExceeded }
	usecase.EnableDynamicTile(jenkinsApi.JenkinsMultiBranchTileType, coreModels.DefaultVariant, &jenkinsModels.MultiBranchParams{}, mockBuilder)

	config, err := readConfig(input)
	if assert.NoError(t, err) {
		usecase.Hydrate(config)
		assert.Len(t, config.Errors, 0)
		assert.Equal(t, jenkinsApi.JenkinsBuildTileType, config.Config.Tiles[0].Type)
		assert.Equal(t, "/jenkins/default?job=test", config.Config.Tiles[0].URL)
	}
}
