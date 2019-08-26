package usecase

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/monitorable/config/models"

	"github.com/monitoror/monitoror/monitorable/config/repository"

	"github.com/stretchr/testify/assert"
)

func initTile(t *testing.T, input string) (tiles *models.Tile) {
	tiles = &models.Tile{}

	err := json.Unmarshal([]byte(input), &tiles)
	assert.NoError(t, err)

	return
}

func TestUsecase_Verify_Success(t *testing.T) {
	input := `
{
	"version" : 3,
  "columns": 4,
  "tiles": [
		{ "type": "EMPTY" }
  ]
}
`
	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.GetConfig(reader)

	if assert.NoError(t, err) {
		useCase := initConfigUsecase()

		err = useCase.Verify(config)
		assert.NoError(t, err)
	}
}

func TestUsecase_Verify_Failed(t *testing.T) {
	input := `
{}
`
	reader := ioutil.NopCloser(strings.NewReader(input))
	config, err := repository.GetConfig(reader)

	if assert.NoError(t, err) {
		useCase := initConfigUsecase()
		err := useCase.Verify(config)

		if assert.Error(t, err) {
			configError := err.(*models.ConfigError)

			assert.Equal(t, 3, configError.Count())
			assert.Contains(t, configError.Error(), `Unsupported "version" field. Must be`)
			assert.Contains(t, configError.Error(), `Missing or invalid "columns" field. Must be a positive integer.`)
			assert.Contains(t, configError.Error(), `Missing or invalid "tiles" field. Must be an array not empty.`)
		}
	}
}

func TestUsecase_VerifyTile_Success(t *testing.T) {
	input := `{ "type": "PORT", "columnSpan": 2, "rowSpan": 2, "params": { "hostname": "bserver.com", "port": 22 } }`

	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 0, configError.Count())
}

func TestUsecase_VerifyTile_Success_Empty(t *testing.T) {
	input := `{ "type": "EMPTY" }`

	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 0, configError.Count())
}

func TestUsecase_VerifyTile_Failed_ParamsInGroup(t *testing.T) {
	input := `
      { "type": "GROUP", "label": "...", "params": {"test": "test"}}
`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unauthorized "params" key in GROUP tile definition.`)
}

func TestUsecase_VerifyTile_Failed_EmptyInGroup(t *testing.T) {
	input := `
      { "type": "GROUP", "label": "...", "tiles": [
          { "type": "EMPTY" }
			]}
`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unauthorized "EMPTY" type in GROUP tile.`)
}

func TestUsecase_VerifyTile_Failed_MissingParamsKey(t *testing.T) {
	input := `{ "type": "PING", "label": "..." }`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Missing "params" key in PING tile definition.`)
}

func TestUsecase_VerifyTile_Success_Group(t *testing.T) {
	input := `
      { "type": "GROUP", "label": "...", "tiles": [
          { "type": "PING", "params": { "hostname": "aserver.com" } },
          { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } }
			]}
`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 0, configError.Count())
}

func TestUsecase_VerifyTile_Failed_GroupInGroup(t *testing.T) {
	input := `
      { "type": "GROUP", "label": "...", "tiles": [
          { "type": "GROUP" }
			]}
`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unauthorized "GROUP" type in GROUP tile.`)
}

func TestUsecase_VerifyTile_Failed_GroupWithWrongTiles(t *testing.T) {
	input := `
      { "type": "GROUP", "label": "..."}
`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Missing or empty "tiles" key in GROUP tile definition.`)
}

func TestUsecase_VerifyTile_Failed_WrongTileType(t *testing.T) {
	input := `{ "type": "PONG", "params": { "hostname": "server.com" } }`

	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unknown "PONG" type in tile definition. Must be`)
}

func TestUsecase_VerifyTile_Failed_InvalidParams(t *testing.T) {
	input := `{ "type": "PING", "params": { "host": "server.com" } }`

	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Invalid params definition for "PING": "{"host":"server.com"}".`)
}

func TestUsecase_VerifyTile_Failed_InvalidColumnSpan(t *testing.T) {
	input := `{ "type": "PING", "columnSpan": -1, "params": { "host": "server.com" } }`

	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Invalid "columnSpan" field. Must be a positive integer.`)
}

func TestUsecase_VerifyTile_Failed_InvalidRowSpan(t *testing.T) {
	input := `{ "type": "PING", "rowSpan": -1, "params": { "host": "server.com" } }`

	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Invalid "rowSpan" field. Must be a positive integer.`)
}

func TestUsecase_VerifyTile_Failed_WrongVariant(t *testing.T) {
	input := `{ "type": "JENKINS-BUILD", "configVariant": "test", "params": { "job": "job1" } }`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unknown "test" variant for JENKINS-BUILD type in tile definition. Must be`)
}

func TestUsecase_VerifyTile_WithDynamicTile(t *testing.T) {
	input := `{ "type": "JENKINS-MULTIBRANCH", "configVariant": "default", "params": { "job": "job1" } }`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 0, configError.Count())
}

func TestUsecase_VerifyTile_WithDynamicTile_WithWrongVariant(t *testing.T) {
	input := `{ "type": "JENKINS-MULTIBRANCH", "configVariant": "test", "params": { "job": "job1" } }`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unknown "test" variant for JENKINS-MULTIBRANCH dynamic type in tile definition. Must be`)
}
