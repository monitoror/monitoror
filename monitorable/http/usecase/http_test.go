package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/monitoror/monitoror/monitorable/http"
	"gopkg.in/yaml.v2"

	"github.com/monitoror/monitoror/models/tiles"

	. "github.com/stretchr/testify/mock"

	"github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/monitorable/http/mocks"
	"github.com/monitoror/monitoror/monitorable/http/models"

	"github.com/stretchr/testify/assert"
)

func TestHttpAny_WithError(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(nil, context.DeadlineExceeded)
	tu := NewHttpUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

	tile, err := tu.HttpAny(&models.HttpAnyParams{Url: "toto"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHtmlAll_WithoutErrors(t *testing.T) {
	for _, testcase := range []struct {
		body            string
		usecaseFunc     func(usecase http.Usecase) (*tiles.HealthTile, error)
		expectedStatus  tiles.TileStatus
		expectedLabel   string
		expectedMessage string
	}{
		{
			// Http Any
			usecaseFunc: func(usecase http.Usecase) (tile *tiles.HealthTile, e error) {
				return usecase.HttpAny(&models.HttpAnyParams{Url: "toto"})
			},
			expectedStatus: tiles.SuccessStatus, expectedLabel: "toto",
		},
		{
			// Http Any with wrong status
			usecaseFunc: func(usecase http.Usecase) (tile *tiles.HealthTile, e error) {
				return usecase.HttpAny(&models.HttpAnyParams{Url: "toto", StatusCodeMin: pointer.ToInt(400), StatusCodeMax: pointer.ToInt(499)})
			},
			expectedStatus: tiles.FailedStatus, expectedLabel: "toto", expectedMessage: "status code 200",
		},
		{
			// Http Raw with matched regex
			body: "errors: 28",
			usecaseFunc: func(usecase http.Usecase) (tile *tiles.HealthTile, e error) {
				return usecase.HttpRaw(&models.HttpRawParams{Url: "toto", Regex: `errors: (\d*)`})
			},
			expectedStatus: tiles.SuccessStatus, expectedLabel: "toto", expectedMessage: "28",
		},
		{
			// Http Raw without matched regex
			body: "api call: 20",
			usecaseFunc: func(usecase http.Usecase) (tile *tiles.HealthTile, e error) {
				return usecase.HttpRaw(&models.HttpRawParams{Url: "toto", Regex: `errors: (\d*)`})
			},
			expectedStatus: tiles.FailedStatus, expectedLabel: "toto", expectedMessage: `pattern not found "errors: (\d*)"`,
		},
		{
			// Http Json
			body: `{"key": "value"}`,
			usecaseFunc: func(usecase http.Usecase) (tile *tiles.HealthTile, e error) {
				return usecase.HttpJson(&models.HttpJsonParams{Url: "toto", Key: ".key"})
			},
			expectedStatus: tiles.SuccessStatus, expectedLabel: "toto", expectedMessage: "value",
		},
		{
			// Http Json missing key
			body: `{"key": "value"}`,
			usecaseFunc: func(usecase http.Usecase) (tile *tiles.HealthTile, e error) {
				return usecase.HttpJson(&models.HttpJsonParams{Url: "toto", Key: ".key2"})
			},
			expectedStatus: tiles.FailedStatus, expectedLabel: "toto", expectedMessage: `unable to lookup for key ".key2"`,
		},
		{
			// Http Json unable to unmarshal
			body: `{"key": "value`,
			usecaseFunc: func(usecase http.Usecase) (tile *tiles.HealthTile, e error) {
				return usecase.HttpYaml(&models.HttpYamlParams{Url: "toto", Key: ".key"})
			},
			expectedStatus: tiles.FailedStatus, expectedLabel: "toto", expectedMessage: `unable to unmarshal content`,
		},
		{
			// Http Yaml
			body: "key: value",
			usecaseFunc: func(usecase http.Usecase) (tile *tiles.HealthTile, e error) {
				return usecase.HttpYaml(&models.HttpYamlParams{Url: "toto", Key: ".key"})
			},
			expectedStatus: tiles.SuccessStatus, expectedLabel: "toto", expectedMessage: "value",
		},
	} {
		mockRepository := new(mocks.Repository)
		mockRepository.On("Get", AnythingOfType("string")).
			Return(&models.Response{StatusCode: 200, Body: []byte(testcase.body)}, nil)
		tu := NewHttpUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

		tile, err := testcase.usecaseFunc(tu)
		if assert.NoError(t, err) {
			assert.Equal(t, testcase.expectedStatus, tile.Status)
			assert.Equal(t, testcase.expectedLabel, tile.Label)
			assert.Equal(t, testcase.expectedMessage, tile.Message)
			mockRepository.AssertNumberOfCalls(t, "Get", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestHttpAny_WithCache(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).
		Return(&models.Response{StatusCode: 200, Body: []byte("test with cache")}, nil)

	tu := NewHttpUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

	tile, err := tu.HttpRaw(&models.HttpRawParams{Url: "toto"})
	if assert.NoError(t, err) {
		assert.Equal(t, "toto", tile.Label)
		assert.Equal(t, "test with cache", tile.Message)
	}

	tile, err = tu.HttpRaw(&models.HttpRawParams{Url: "toto"})
	if assert.NoError(t, err) {
		assert.Equal(t, "toto", tile.Label)
		assert.Equal(t, "test with cache", tile.Message)
	}
	mockRepository.AssertNumberOfCalls(t, "Get", 1)
	mockRepository.AssertExpectations(t)
}

func TestHttpUsecase_CheckStatusCode(t *testing.T) {
	httpAny := &models.HttpAnyParams{}
	assert.True(t, checkStatusCode(httpAny, 301))
	assert.False(t, checkStatusCode(httpAny, 404))

	httpAny.StatusCodeMin = pointer.ToInt(200)
	httpAny.StatusCodeMax = pointer.ToInt(399)
	assert.True(t, checkStatusCode(httpAny, 301))
	assert.False(t, checkStatusCode(httpAny, 404))
}

func TestHttpUsecase_Match(t *testing.T) {
	httpRaw := &models.HttpRawParams{}
	match, substring := matchRegex(httpRaw, "test")
	assert.True(t, match)
	assert.Equal(t, "test", substring)

	httpRaw.Regex = "test"
	match, substring = matchRegex(httpRaw, "test 2")
	assert.True(t, match)
	assert.Equal(t, "test 2", substring)

	httpRaw.Regex = `test (\d)`
	match, substring = matchRegex(httpRaw, "test 2")
	assert.True(t, match)
	assert.Equal(t, "2", substring)

	httpRaw.Regex = `toto (\d)`
	match, substring = matchRegex(httpRaw, "test 2")
	assert.False(t, match)
	assert.Equal(t, "", substring)
}

func TestHttpUsecase_LookupKey_Json(t *testing.T) {
	input := `
{
	"bloc1": {
		"bloc.2": [
			{ "value": "YEAH !!" },
			{ "value": "NOOO !!" },
			{ "value": "NOOO !!" }
		]
	}
}
`
	httpJson := &models.HttpJsonParams{}
	httpJson.Key = `.bloc1."bloc.2".[0].value`

	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	if assert.NoError(t, err) {
		found, value := lookupKey(httpJson, data)
		assert.True(t, found)
		assert.Equal(t, "YEAH !!", value)
	}
}

func TestHttpUsecase_LookupKey_Yaml(t *testing.T) {
	input := `
bloc1: 
  bloc.2: 
    - name: test1
      value: "YEAH !!"
    - name: test2
      value: "NOOO !!"
`
	httpYaml := &models.HttpYamlParams{}
	httpYaml.Key = `.bloc1."bloc.2".[0].value`

	var data interface{}
	err := yaml.Unmarshal([]byte(input), &data)
	if assert.NoError(t, err) {
		found, value := lookupKey(httpYaml, data)
		assert.True(t, found)
		assert.Equal(t, "YEAH !!", value)
	}
}

func TestHttpUsecase_LookupKey_MissingKey(t *testing.T) {
	input := `
bloc1: 
  bloc.2: 
    - name: test1
      value: "YEAH !!"
    - name: test2
      value: "NOOO !!"
`
	httpYaml := &models.HttpYamlParams{}
	httpYaml.Key = `.bloc1."bloc.2".[0].value2`

	var data interface{}
	err := yaml.Unmarshal([]byte(input), &data)
	if assert.NoError(t, err) {
		found, _ := lookupKey(httpYaml, data)
		assert.False(t, found)
	}
}
