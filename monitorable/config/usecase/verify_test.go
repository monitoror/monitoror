package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/monitorable/config/models"

	"github.com/monitoror/monitoror/monitorable/config/repository"

	"github.com/stretchr/testify/assert"
)

func initTile(t *testing.T, input string) (tiles map[string]interface{}) {
	tiles = make(map[string]interface{})

	err := json.Unmarshal([]byte(input), &tiles)
	assert.NoError(t, err)

	return
}

func TestUsecase_Verify_Success(t *testing.T) {
	input := `
{
	"version" : 1,
  "columns": 4,
  "tiles": [
		{ "type": "empty" }
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
			assert.Contains(t, configError.Error(), fmt.Sprintf(`Unsupported "version" field. Must be %s.`, keys(SupportedVersions)))
			assert.Contains(t, configError.Error(), `Missing or invalid "columns" field. Must be a positive integer.`)
			assert.Contains(t, configError.Error(), `Missing or invalid "tiles" field. Must be an array not empty.`)
		}
	}
}

func TestUsecase_VerifyTile_Success(t *testing.T) {
	input := `{ "type": "port", "params": { "hostname": "bserver.com", "port": 22 } }`

	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 0, configError.Count())
}

func TestUsecase_VerifyTile_Success_Empty(t *testing.T) {
	input := `{ "type": "empty" }`

	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 0, configError.Count())
}

func TestUsecase_VerifyTile_Failed_WrongKey(t *testing.T) {
	input := `{ "test": "empty" }`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unknown key "test" in tile definition. Must be`)
}

func TestUsecase_VerifyTile_Failed_ParamsInGroup(t *testing.T) {
	input := `
      { "type": "group", "label": "...", "params": [
          { "type": "ping", "params": { "hostname": "aserver.com" } }
			]}
`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unauthorized "params" key in group tile definition.`)
}

func TestUsecase_VerifyTile_Failed_EmptyInGroup(t *testing.T) {
	input := `
      { "type": "group", "label": "...", "tiles": [
          { "type": "empty" }
			]}
`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unauthorized "empty" type in group tile.`)
}

func TestUsecase_VerifyTile_Failed_MissingParamsKey(t *testing.T) {
	input := `{ "type": "ping", "label": "..." }`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Missing "params" key in ping tile definition.`)
}

func TestUsecase_VerifyTile_Success_Group(t *testing.T) {
	input := `
      { "type": "group", "label": "...", "tiles": [
          { "type": "ping", "params": { "hostname": "aserver.com" } },
          { "type": "port", "params": { "hostname": "bserver.com", "port": 22 } }
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
      { "type": "group", "label": "...", "tiles": [
          { "type": "group", "params": "test" }
			]}
`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unauthorized "group" type in group tile.`)
}

func TestUsecase_VerifyTile_Failed_GroupWithWrongTiles(t *testing.T) {
	input := `
      { "type": "group", "label": "...", "tiles": "test"}
`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Incorrect "tiles" key in group tile definition.`)
}

func TestUsecase_VerifyTile_Failed_GroupWithWrongTile(t *testing.T) {
	input := `
      { "type": "group", "label": "...", "tiles": [
          { "type": "ping", "params": { "hostname": "aserver.com" } },
          "test"
			]}
`
	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Incorrect array element "test" in group definition.`)
}

func TestUsecase_VerifyTile_Failed_WrongTileType(t *testing.T) {
	input := `{ "type": "pong", "params": { "hostname": "server.com" } }`

	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Unknown "pong" type in tile definition. Must be`)
}

func TestUsecase_VerifyTile_Failed_InvalidParams(t *testing.T) {
	input := `{ "type": "ping", "params": { "host": "server.com" } }`

	configError := &models.ConfigError{}

	tile := initTile(t, input)
	useCase := initConfigUsecase()

	useCase.verifyTile(tile, false, configError)

	assert.Equal(t, 1, configError.Count())
	assert.Contains(t, configError.Error(), `Invalid params definition for "ping": "{"host":"server.com"}".`)
}
