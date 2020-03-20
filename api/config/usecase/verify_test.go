package usecase

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/config"
	jenkinsApi "github.com/monitoror/monitoror/monitorables/jenkins/api"
	jenkinsModels "github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"

	"github.com/stretchr/testify/assert"
)

func initTile(t *testing.T, rawConfig string) (tiles *models.TileConfig) {
	tiles = &models.TileConfig{}

	err := json.Unmarshal([]byte(rawConfig), &tiles)
	assert.NoError(t, err)

	return
}

func TestUsecase_Verify_Success(t *testing.T) {
	rawConfig := fmt.Sprintf(`
{
	"version" : %q,
  "columns": 4,
  "tiles": [
		{ "type": "EMPTY" }
  ]
}
`, CurrentVersion)

	conf, err := readConfig(rawConfig)
	if assert.NoError(t, err) {
		usecase := initConfigUsecase(nil, nil)
		usecase.Verify(conf)

		assert.Len(t, conf.Errors, 0)
	}
}

func TestUsecase_Verify_SuccessWithOptionalParameters(t *testing.T) {
	rawConfig := fmt.Sprintf(`
{
	"version" : %q,
  "columns": 4,
  "zoom": 2.5,
  "tiles": [
		{ "type": "EMPTY" }
  ]
}
`, CurrentVersion)

	conf, err := readConfig(rawConfig)

	if assert.NoError(t, err) {
		usecase := initConfigUsecase(nil, nil)
		usecase.Verify(conf)

		assert.Len(t, conf.Errors, 0)
	}
}

func TestUsecase_Verify_Failed(t *testing.T) {
	for _, testcase := range []struct {
		rawConfig string
		errorID   models.ConfigErrorID
		errorData models.ConfigErrorData
	}{
		{
			rawConfig: `{}`,
			errorID:   models.ConfigErrorMissingRequiredField,
			errorData: models.ConfigErrorData{FieldName: "version"},
		},
		{
			rawConfig: `{"version": "0.0"}`,
			errorID:   models.ConfigErrorUnsupportedVersion,
			errorData: models.ConfigErrorData{
				FieldName: "version",
				Value:     `"0.0"`,
				Expected:  fmt.Sprintf(`%q >= version >= %q`, MinimalVersion, CurrentVersion),
			},
		},
		{
			rawConfig: `{"version": "999.999"}`, // Don't let me go that far ^^'
			errorID:   models.ConfigErrorUnsupportedVersion,
			errorData: models.ConfigErrorData{
				FieldName: "version",
				Value:     `"999.999"`,
				Expected:  fmt.Sprintf(`%q >= version >= %q`, MinimalVersion, CurrentVersion),
			},
		},
		{
			rawConfig: fmt.Sprintf(`{"version": %q, "tiles": [{ "type": "EMPTY" }]}`, CurrentVersion),
			errorID:   models.ConfigErrorMissingRequiredField,
			errorData: models.ConfigErrorData{
				FieldName: "columns",
			},
		},
		{
			rawConfig: fmt.Sprintf(`{"version": %q, "columns": 0, "tiles": [{ "type": "EMPTY" }]}`, CurrentVersion),
			errorID:   models.ConfigErrorInvalidFieldValue,
			errorData: models.ConfigErrorData{
				FieldName: "columns",
				Value:     `0`,
				Expected:  "columns > 0",
			},
		},
		{
			rawConfig: fmt.Sprintf(`{"version": %q, "columns": 1, "zoom": 0, "tiles": [{ "type": "EMPTY" }]}`, CurrentVersion),
			errorID:   models.ConfigErrorInvalidFieldValue,
			errorData: models.ConfigErrorData{
				FieldName: "zoom",
				Value:     `0`,
				Expected:  "0 < zoom <= 10",
			},
		},
		{
			rawConfig: fmt.Sprintf(`{"version": %q, "columns": 1, "zoom": 20, "tiles": [{ "type": "EMPTY" }]}`, CurrentVersion),
			errorID:   models.ConfigErrorInvalidFieldValue,
			errorData: models.ConfigErrorData{
				FieldName: "zoom",
				Value:     `20`,
				Expected:  "0 < zoom <= 10",
			},
		},
		{
			rawConfig: fmt.Sprintf(`{"version": %q, "columns": 1}`, CurrentVersion),
			errorID:   models.ConfigErrorMissingRequiredField,
			errorData: models.ConfigErrorData{
				FieldName: "tiles",
			},
		},
		{
			rawConfig: fmt.Sprintf(`{"version": %q, "columns": 1, "tiles": []}`, CurrentVersion),
			errorID:   models.ConfigErrorInvalidFieldValue,
			errorData: models.ConfigErrorData{
				FieldName:     "tiles",
				ConfigExtract: fmt.Sprintf(`{"version":%q,"columns":1,"tiles":[]}`, CurrentVersion),
			},
		},
	} {
		conf, err := readConfig(testcase.rawConfig)
		if assert.NoError(t, err) {
			usecase := initConfigUsecase(nil, nil)
			usecase.Verify(conf)
			if assert.Len(t, conf.Errors, 1) {
				assert.Equal(t, testcase.errorID, conf.Errors[0].ID)
				assert.Equal(t, testcase.errorData, conf.Errors[0].Data)
			}
		}
	}
}

func TestUsecase_VerifyTile_Success(t *testing.T) {
	rawConfig := `{ "type": "PORT", "columnSpan": 2, "rowSpan": 2, "params": { "hostname": "bserver.com", "port": 22 } }`
	conf := &models.ConfigBag{}

	tile := initTile(t, rawConfig)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, nil)

	assert.Len(t, conf.Errors, 0)
}

func TestUsecase_VerifyTile_Success_Empty(t *testing.T) {
	rawConfig := `{ "type": "EMPTY" }`
	conf := &models.ConfigBag{}

	tile := initTile(t, rawConfig)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, nil)

	assert.Len(t, conf.Errors, 0)
}

func TestUsecase_VerifyTile_Success_Group(t *testing.T) {
	rawConfig := `
      { "type": "GROUP", "label": "...", "tiles": [
          { "type": "PING", "params": { "hostname": "aserver.com" } },
          { "type": "PORT", "params": { "hostname": "bserver.com", "port": 22 } }
			]}
`
	conf := &models.ConfigBag{}

	tile := initTile(t, rawConfig)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, nil)

	assert.Len(t, conf.Errors, 0)
}

func TestUsecase_VerifyTile_Failed(t *testing.T) {
	for _, testcase := range []struct {
		rawConfig string
		errorID   models.ConfigErrorID
		errorData models.ConfigErrorData
	}{
		{
			rawConfig: `{ "type": "PING", "columnSpan": -1, "params": { "hostname": "server.com" } }`,
			errorID:   models.ConfigErrorInvalidFieldValue,
			errorData: models.ConfigErrorData{
				FieldName: "columnSpan",
				Value:     "-1",
				Expected:  "columnSpan > 0",
			},
		},
		{
			rawConfig: `{ "type": "PING", "rowSpan": -1, "params": { "hostname": "server.com" } }`,
			errorID:   models.ConfigErrorInvalidFieldValue,
			errorData: models.ConfigErrorData{
				FieldName: "rowSpan",
				Value:     "-1",
				Expected:  "rowSpan > 0",
			},
		},
		{
			rawConfig: `
					{ "type": "GROUP", "tiles": [
							{ "type": "EMPTY" }
					]}
		`,
			errorID: models.ConfigErrorUnauthorizedSubtileType,
			errorData: models.ConfigErrorData{
				ConfigExtract:          `{"type":"GROUP","tiles":[{"type":"EMPTY"}]}`,
				ConfigExtractHighlight: `{"type":"EMPTY"}`,
			},
		},
		{
			rawConfig: `
					{ "type": "GROUP", "tiles": [
							{ "type": "GROUP" }
					]}
		`,
			errorID: models.ConfigErrorUnauthorizedSubtileType,
			errorData: models.ConfigErrorData{
				ConfigExtract:          `{"type":"GROUP","tiles":[{"type":"GROUP"}]}`,
				ConfigExtractHighlight: `{"type":"GROUP"}`,
			},
		},
		{
			rawConfig: `{ "type": "GROUP", "params": {"test": "test"}}`,
			errorID:   models.ConfigErrorUnauthorizedField,
			errorData: models.ConfigErrorData{
				FieldName:     "params",
				ConfigExtract: `{"type":"GROUP","params":{"test":"test"}}`,
			},
		},
		{
			rawConfig: `{ "type": "GROUP"}`,
			errorID:   models.ConfigErrorMissingRequiredField,
			errorData: models.ConfigErrorData{
				FieldName:     "tiles",
				ConfigExtract: `{"type":"GROUP"}`,
			},
		},
		{
			rawConfig: `{ "type": "GROUP", "tiles": []}`,
			errorID:   models.ConfigErrorInvalidFieldValue,
			errorData: models.ConfigErrorData{
				FieldName:     "tiles",
				ConfigExtract: `{"type":"GROUP"}`,
			},
		},
		{
			rawConfig: `{ "type": "PING" }`,
			errorID:   models.ConfigErrorMissingRequiredField,
			errorData: models.ConfigErrorData{
				FieldName:     "params",
				ConfigExtract: `{"type":"PING"}`,
			},
		},
		{
			rawConfig: `{ "type": "PING", "params": { } }`,
			errorID:   models.ConfigErrorInvalidFieldValue,
			errorData: models.ConfigErrorData{
				FieldName:     "params",
				ConfigExtract: `{"type":"PING","configVariant":"default"}`,
			},
		},
		{
			rawConfig: `{ "type": "PING", "params": { "host": "server.com" } }`,
			errorID:   models.ConfigErrorUnknownField,
			errorData: models.ConfigErrorData{
				FieldName:     "host",
				ConfigExtract: `{"type":"PING","params":{"host":"server.com"},"configVariant":"default"}`,
				Expected:      "hostname",
			},
		},
		{
			rawConfig: `{ "type": "JENKINS-BUILD", "configVariant": "test", "params": { "job": "job1" } }`,
			errorID:   models.ConfigErrorUnknownVariant,
			errorData: models.ConfigErrorData{
				FieldName:     "configVariant",
				Value:         `"test"`,
				ConfigExtract: `{"type":"JENKINS-BUILD","params":{"job":"job1"},"configVariant":"test"}`,
			},
		},
	} {
		conf := &models.ConfigBag{}
		tile := initTile(t, testcase.rawConfig)
		usecase := initConfigUsecase(nil, nil)
		usecase.verifyTile(conf, tile, nil)

		if assert.Len(t, conf.Errors, 1) {
			assert.Equal(t, testcase.errorID, conf.Errors[0].ID)
			assert.Equal(t, testcase.errorData, conf.Errors[0].Data)
		}
	}
}

func TestUsecase_VerifyTile_Failed_WrongTileType(t *testing.T) {
	rawConfig := `{ "type": "PONG", "params": { "hostname": "server.com" } }`

	conf := &models.ConfigBag{}
	tile := initTile(t, rawConfig)
	usecase := initConfigUsecase(nil, nil)
	usecase.verifyTile(conf, tile, nil)

	if assert.Len(t, conf.Errors, 1) {
		assert.Equal(t, models.ConfigErrorUnknownTileType, conf.Errors[0].ID)
		assert.Equal(t, "type", conf.Errors[0].Data.FieldName)
		assert.Equal(t, `{"type":"PONG","params":{"hostname":"server.com"}}`, conf.Errors[0].Data.ConfigExtract)
	}
}

func TestUsecase_VerifyTile_WithDynamicTile(t *testing.T) {
	rawConfig := `{ "type": "JENKINS-MULTIBRANCH", "configVariant": "default", "params": { "job": "job1" } }`
	conf := &models.ConfigBag{}

	tile := initTile(t, rawConfig)

	params := make(map[string]interface{})
	params["job"] = "test"
	mockBuilder := func(_ interface{}) ([]builder.Result, error) {
		return []builder.Result{{TileType: jenkinsApi.JenkinsBuildTileType, Params: params}}, nil
	}

	usecase := initConfigUsecase(nil, nil)
	usecase.EnableDynamicTile(jenkinsApi.JenkinsMultiBranchTileType, config.DefaultVariant, &jenkinsModels.MultiBranchParams{}, mockBuilder)
	usecase.verifyTile(conf, tile, nil)

	assert.Len(t, conf.Errors, 0)
}

func TestUsecase_VerifyTile_WithDynamicTile_WithWrongVariant(t *testing.T) {
	rawConfig := `{ "type": "JENKINS-MULTIBRANCH", "configVariant": "test", "params": { "job": "job1" } }`
	conf := &models.ConfigBag{}

	tile := initTile(t, rawConfig)

	params := make(map[string]interface{})
	params["job"] = "test"
	mockBuilder := func(_ interface{}) ([]builder.Result, error) {
		return []builder.Result{{TileType: jenkinsApi.JenkinsBuildTileType, Params: params}}, nil
	}

	usecase := initConfigUsecase(nil, nil)
	usecase.EnableDynamicTile(jenkinsApi.JenkinsMultiBranchTileType, config.DefaultVariant, &jenkinsModels.MultiBranchParams{}, mockBuilder)
	usecase.verifyTile(conf, tile, nil)

	if assert.Len(t, conf.Errors, 1) {
		assert.Equal(t, models.ConfigErrorUnknownVariant, conf.Errors[0].ID)
		assert.Equal(t, "configVariant", conf.Errors[0].Data.FieldName)
		assert.Equal(t, `"test"`, conf.Errors[0].Data.Value)
		assert.Equal(t, `{"type":"JENKINS-MULTIBRANCH","params":{"job":"job1"},"configVariant":"test"}`, conf.Errors[0].Data.ConfigExtract)
	}
}
