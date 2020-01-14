package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/monitorable/config/repository"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	_jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	. "github.com/monitoror/monitoror/pkg/monitoror/builder/mocks"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initTile(t *testing.T, input string) (tiles *models.Tile) {
	tiles = &models.Tile{}

	err := json.Unmarshal([]byte(input), &tiles)
	assert.NoError(t, err)

	return
}

func TestUsecase_Verify_Success(t *testing.T) {
	input := fmt.Sprintf(`
{
	"version" : %d,
  "columns": 4,
  "tiles": [
		{ "type": "EMPTY" }
  ]
}
`, CurrentVersion)

	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.ReadConfig(reader)

	if assert.NoError(t, err) {
		usecase := initConfigUsecase(nil, nil)
		usecase.Verify(config)
		assert.Len(t, config.Errors, 0)
	}
}

func TestUsecase_Verify_MissingVersion(t *testing.T) {
	input := `{}`
	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.ReadConfig(reader)

	if assert.NoError(t, err) {
		usecase := initConfigUsecase(nil, nil)
		usecase.Verify(config)
		if assert.Len(t, config.Errors, 1) {
			assert.Contains(t, config.Errors[0], `Missing "version" field. Must be`)
		}
	}
}

func TestUsecase_Verify_UnknownVersion(t *testing.T) {
	input := `
{"version": 0}
`
	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.ReadConfig(reader)

	if assert.NoError(t, err) {
		usecase := initConfigUsecase(nil, nil)
		usecase.Verify(config)
		if assert.Len(t, config.Errors, 1) {
			assert.Contains(t, config.Errors[0], `Unsupported "version" field. Must be`)
		}
	}
}

func TestUsecase_Verify_Failed(t *testing.T) {
	input := fmt.Sprintf(`
{"version": %d}
`, CurrentVersion)

	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.ReadConfig(reader)

	if assert.NoError(t, err) {
		usecase := initConfigUsecase(nil, nil)
		usecase.Verify(config)
		if assert.Len(t, config.Errors, 2) {
			assert.Contains(t, config.Errors, `Missing or invalid "columns" field. Must be a positive integer.`)
			assert.Contains(t, config.Errors, `Missing or invalid "tiles" field. Must be an array not empty.`)
		}
	}
}

func TestUsecase_VerifyTile_Success(t *testing.T) {
	input := `{ "type": "PORT", "columnSpan": 2, "rowSpan": 2, "params": { "hostname": "bserver.com", "port": 22 } }`
	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 0)
}

func TestUsecase_VerifyTile_Success_Empty(t *testing.T) {
	input := `{ "type": "EMPTY" }`
	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 0)
}

func TestUsecase_VerifyTile_Failed_ParamsInGroup(t *testing.T) {
	input := `{ "type": "GROUP", "label": "...", "params": {"test": "test"}}`
	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Unauthorized "params" key in GROUP tile definition.`)
}

func TestUsecase_VerifyTile_Failed_EmptyInGroup(t *testing.T) {
	input := `
      { "type": "GROUP", "label": "...", "tiles": [
          { "type": "EMPTY" }
			]}
`
	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Unauthorized "EMPTY" type in GROUP tile.`)
}

func TestUsecase_VerifyTile_Failed_MissingParamsKey(t *testing.T) {
	input := `{ "type": "PING", "label": "..." }`
	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Missing "params" key in PING tile definition.`)
}

func TestUsecase_VerifyTile_Success_Group(t *testing.T) {
	input := `
      { "type": "GROUP", "label": "...", "tiles": [
          { "type": "PING", "params": { "hostname": "aserver.com" } },
          { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } }
			]}
`
	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 0)
}

func TestUsecase_VerifyTile_Failed_GroupInGroup(t *testing.T) {
	input := `
      { "type": "GROUP", "label": "...", "tiles": [
          { "type": "GROUP" }
			]}
`
	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Unauthorized "GROUP" type in GROUP tile.`)
}

func TestUsecase_VerifyTile_Failed_GroupWithWrongTiles(t *testing.T) {
	input := `
     { "type": "GROUP", "label": "..."}
`
	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Missing or empty "tiles" key in GROUP tile definition.`)
}

func TestUsecase_VerifyTile_Failed_WrongTileType(t *testing.T) {
	input := `{ "type": "PONG", "params": { "hostname": "server.com" } }`

	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Unknown "PONG" type in tile definition. Must be`)
}

func TestUsecase_VerifyTile_Failed_InvalidParams(t *testing.T) {
	input := `{ "type": "PING", "params": { "host": "server.com" } }`

	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Invalid params definition for "PING": "{"host":"server.com"}".`)
}

func TestUsecase_VerifyTile_Failed_InvalidColumnSpan(t *testing.T) {
	input := `{ "type": "PING", "columnSpan": -1, "params": { "host": "server.com" } }`

	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Invalid "columnSpan" field. Must be a positive integer.`)
}

func TestUsecase_VerifyTile_Failed_InvalidRowSpan(t *testing.T) {
	input := `{ "type": "PING", "rowSpan": -1, "params": { "host": "server.com" } }`

	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Invalid "rowSpan" field. Must be a positive integer.`)
}

func TestUsecase_VerifyTile_Failed_WrongVariant(t *testing.T) {
	input := `{ "type": "JENKINS-BUILD", "configVariant": "test", "params": { "job": "job1" } }`
	conf := &models.Config{}

	tile := initTile(t, input)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Unknown "test" variant for JENKINS-BUILD type in tile definition. Must be`)
}

func TestUsecase_VerifyTile_WithDynamicTile(t *testing.T) {
	input := `{ "type": "JENKINS-MULTIBRANCH", "configVariant": "default", "params": { "job": "job1" } }`
	conf := &models.Config{}

	tile := initTile(t, input)

	params := make(map[string]interface{})
	params["job"] = "test"
	mockBuilder := new(DynamicTileBuilder)
	mockBuilder.On("ListDynamicTile", Anything).Return([]builder.Result{{TileType: jenkins.JenkinsBuildTileType, Params: params}}, nil)

	usecase := initConfigUsecase(nil, nil)
	usecase.RegisterDynamicTile(jenkins.JenkinsMultiBranchTileType, &_jenkinsModels.MultiBranchParams{}, mockBuilder)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 0)
}

func TestUsecase_VerifyTile_WithDynamicTile_WithWrongVariant(t *testing.T) {
	input := `{ "type": "JENKINS-MULTIBRANCH", "configVariant": "test", "params": { "job": "job1" } }`
	conf := &models.Config{}

	tile := initTile(t, input)

	params := make(map[string]interface{})
	params["job"] = "test"
	mockBuilder := new(DynamicTileBuilder)
	mockBuilder.On("ListDynamicTile", Anything).Return([]builder.Result{{TileType: jenkins.JenkinsBuildTileType, Params: params}}, nil)

	usecase := initConfigUsecase(nil, nil)
	usecase.RegisterDynamicTile(jenkins.JenkinsMultiBranchTileType, &_jenkinsModels.MultiBranchParams{}, mockBuilder)
	usecase.verifyTile(conf, tile, false)

	assert.Len(t, conf.Errors, 1)
	assert.Contains(t, conf.Errors[0], `Unknown "test" variant for JENKINS-MULTIBRANCH dynamic type in tile definition. Must be`)
}
