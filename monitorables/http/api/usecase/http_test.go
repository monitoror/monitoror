package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/mocks"
	"github.com/monitoror/monitoror/monitorables/http/api/models"

	"github.com/AlekSi/pointer"
	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v2"
)

func TestHTTPStatus_WithError(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(nil, context.DeadlineExceeded)
	tu := NewHTTPUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

	tile, err := tu.HTTPStatus(&models.HTTPStatusParams{URL: "url"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHtmlAll_WithoutErrors(t *testing.T) {
	for _, testcase := range []struct {
		body                string
		usecaseFunc         func(usecase api.Usecase) (*coreModels.Tile, error)
		expectedStatus      coreModels.TileStatus
		expectedLabel       string
		expectedMessage     string
		expectedValueUnit   coreModels.TileValuesUnit
		expectedValueValues []string
	}{
		{
			// HTTP Status
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPStatus(&models.HTTPStatusParams{URL: "url"})
			},
			expectedStatus: coreModels.SuccessStatus, expectedLabel: "url",
		},
		{
			// HTTP Status with wrong status
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPStatus(&models.HTTPStatusParams{URL: "url", StatusCodeMin: pointer.ToInt(400), StatusCodeMax: pointer.ToInt(499)})
			},
			expectedStatus: coreModels.FailedStatus, expectedLabel: "url", expectedMessage: "status code 200",
		},
		{
			// HTTP Raw with matched regex
			body: "test",
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPStatus(&models.HTTPStatusParams{URL: "url"})
			},
			expectedStatus: coreModels.SuccessStatus, expectedLabel: "url",
		},
		{
			// HTTP Raw with matched regex
			body: "errors: 28",
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPRaw(&models.HTTPRawParams{URL: "url", Regex: `errors: (\d*)`})
			},
			expectedStatus: coreModels.SuccessStatus, expectedLabel: "url", expectedValueUnit: coreModels.NumberUnit, expectedValueValues: []string{"28"},
		},
		{
			// HTTP Raw without matched regex
			body: "api call: 20",
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPRaw(&models.HTTPRawParams{URL: "url", Regex: `errors: (\d*)`})
			},
			expectedStatus: coreModels.FailedStatus, expectedLabel: "url", expectedValueUnit: coreModels.RawUnit, expectedValueValues: []string{`api call: 20`},
		},
		{
			// HTTP Json
			body: `{"key": "value"}`,
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPFormatted(&models.HTTPFormattedParams{URL: "url", Format: models.JSONFormat, Key: "key"})
			},
			expectedStatus: coreModels.SuccessStatus, expectedLabel: "url", expectedValueUnit: coreModels.RawUnit, expectedValueValues: []string{"value"},
		},
		{
			// HTTP Json with key jq like
			body: `{"key": "value"}`,
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPFormatted(&models.HTTPFormattedParams{URL: "url", Format: models.JSONFormat, Key: ".key"})
			},
			expectedStatus: coreModels.SuccessStatus, expectedLabel: "url", expectedValueUnit: coreModels.RawUnit, expectedValueValues: []string{"value"},
		},
		{
			// HTTP Json with long float
			body: `{"key": 123456789 }`,
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPFormatted(&models.HTTPFormattedParams{URL: "url", Format: models.JSONFormat, Key: "key"})
			},
			expectedStatus: coreModels.SuccessStatus, expectedLabel: "url", expectedValueUnit: coreModels.NumberUnit, expectedValueValues: []string{"123456789"},
		},
		{
			// HTTP Json missing key
			body: `{"key": "value"}`,
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPFormatted(&models.HTTPFormattedParams{URL: "url", Format: models.JSONFormat, Key: "key2"})
			},
			expectedStatus: coreModels.FailedStatus, expectedLabel: "url", expectedMessage: `unable to lookup for key "key2"`,
		},
		{
			// HTTP Json unable to unmarshal
			body: `{"key": "value`,
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPFormatted(&models.HTTPFormattedParams{URL: "url", Format: models.JSONFormat, Key: "key"})
			},
			expectedStatus: coreModels.FailedStatus, expectedLabel: "url", expectedMessage: `unable to unmarshal content`,
		},
		{
			// HTTP XML
			body: `<check><status test="2">OK</status></check>`,
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPFormatted(&models.HTTPFormattedParams{URL: "url", Format: models.XMLFormat, Key: "check.status.#content"})
			},
			expectedStatus: coreModels.SuccessStatus, expectedLabel: "url", expectedValueUnit: coreModels.RawUnit, expectedValueValues: []string{"OK"},
		},
		{
			// HTTP XML unable to convert to json
			body: `<check><status test="2">OK</stat`,
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPFormatted(&models.HTTPFormattedParams{URL: "url", Format: models.XMLFormat, Key: "check.status.#content"})
			},
			expectedStatus: coreModels.FailedStatus, expectedLabel: "url", expectedMessage: "unable to convert xml to json",
		},
		{
			// HTTP YAML
			body: "key: value",
			usecaseFunc: func(usecase api.Usecase) (*coreModels.Tile, error) {
				return usecase.HTTPFormatted(&models.HTTPFormattedParams{URL: "url", Format: models.YAMLFormat, Key: "key"})
			},
			expectedStatus: coreModels.SuccessStatus, expectedLabel: "url", expectedValueUnit: coreModels.RawUnit, expectedValueValues: []string{"value"},
		},
	} {
		mockRepository := new(mocks.Repository)
		mockRepository.On("Get", AnythingOfType("string")).
			Return(&models.Response{StatusCode: 200, Body: []byte(testcase.body)}, nil)
		tu := NewHTTPUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

		tile, err := testcase.usecaseFunc(tu)
		if assert.NoError(t, err) {
			assert.Equal(t, testcase.expectedStatus, tile.Status)
			assert.Equal(t, testcase.expectedLabel, tile.Label)
			assert.Equal(t, testcase.expectedMessage, tile.Message)
			if tile.Value != nil {
				assert.Equal(t, testcase.expectedValueUnit, tile.Value.Unit)
				assert.Equal(t, testcase.expectedValueValues, tile.Value.Values)
			}
			mockRepository.AssertNumberOfCalls(t, "Get", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestHTTPStatus_WithCache(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).
		Return(&models.Response{StatusCode: 200, Body: []byte("test with cache")}, nil)

	tu := NewHTTPUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

	tile, err := tu.HTTPRaw(&models.HTTPRawParams{URL: "url"})
	if assert.NoError(t, err) {
		assert.Equal(t, "url", tile.Label)
		assert.Equal(t, "test with cache", tile.Value.Values[0])
	}

	tile, err = tu.HTTPRaw(&models.HTTPRawParams{URL: "url"})
	if assert.NoError(t, err) {
		assert.Equal(t, "url", tile.Label)
		assert.Equal(t, "test with cache", tile.Value.Values[0])
	}
	mockRepository.AssertNumberOfCalls(t, "Get", 1)
	mockRepository.AssertExpectations(t)
}

func TestHTTPProxy_UnreachableUrl(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(nil, context.DeadlineExceeded)
	tu := NewHTTPUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

	tile, err := tu.HTTPProxy(&models.HTTPProxyParams{URL: "url"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.Equal(t, `unable to get "url"`, err.Error())
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHTTPProxy_WrongStatusError(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 404}, nil)
	tu := NewHTTPUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

	tile, err := tu.HTTPProxy(&models.HTTPProxyParams{URL: "url"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.Equal(t, "wrong status code for url", err.Error())
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHTTPProxy_UnmarshallError(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 200, Body: []byte("test")}, nil)
	tu := NewHTTPUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

	tile, err := tu.HTTPProxy(&models.HTTPProxyParams{URL: "url"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		assert.Equal(t, `unable to parse "url" into tile structure`, err.Error())
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHTTPProxy_ContentError(t *testing.T) {
	for _, testcase := range []struct {
		body          string
		expectedError string
	}{
		{
			body:          "{}",
			expectedError: `unauthorized tile.status: ""`,
		},
		{
			body:          `{"status": ""}`,
			expectedError: `unauthorized tile.status: ""`,
		},
		{
			body:          `{"status": "TEST"}`,
			expectedError: `unauthorized tile.status: "TEST"`,
		},
		{
			body:          `{"status": "SUCCESS", "value": {}, "build": {}}`,
			expectedError: `tile.value and tile.build are exclusive`,
		},
		{
			body:          `{"status": "SUCCESS", "value": {"unit": ""}}`,
			expectedError: `unauthorized tile.value.unit: ""`,
		},
		{
			body:          `{"status": "SUCCESS", "value": {"unit": "TEST"}}`,
			expectedError: `unauthorized tile.value.unit: "TEST"`,
		},
		{
			body:          `{"status": "SUCCESS", "value": {"unit": "RAW"}}`,
			expectedError: `unauthorized empty tile.value.values`,
		},
		{
			body:          `{"status": "RUNNING"}`,
			expectedError: `unauthorized empty tile.build with "RUNNING" tile.status`,
		},
		{
			body:          `{"status": "RUNNING", "build": {}}`,
			expectedError: `unauthorized empty tile.build.duration with "RUNNING" tile.status`,
		},
	} {
		mockRepository := new(mocks.Repository)
		mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 200, Body: []byte(testcase.body)}, nil)
		tu := NewHTTPUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

		tile, err := tu.HTTPProxy(&models.HTTPProxyParams{URL: "url"})
		if assert.Error(t, err) {
			assert.Nil(t, tile)
			assert.Equal(t, testcase.expectedError, err.Error())
			mockRepository.AssertNumberOfCalls(t, "Get", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestHTTPProxy_Success(t *testing.T) {
	content := `{"status": "RUNNING", "build": {"duration": 10}}`

	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 200, Body: []byte(content)}, nil)
	tu := NewHTTPUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

	tile, err := tu.HTTPProxy(&models.HTTPProxyParams{URL: "url"})
	if assert.NoError(t, err) {
		assert.Equal(t, api.HTTPProxyTileType, tile.Type)
		assert.Equal(t, pointer.ToInt64(0), tile.Build.EstimatedDuration)
		assert.Equal(t, coreModels.UnknownStatus, tile.Build.PreviousStatus)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHTTPUsecase_CheckStatusCode(t *testing.T) {
	httpAny := &models.HTTPStatusParams{}
	assert.True(t, checkStatusCode(httpAny, 301))
	assert.False(t, checkStatusCode(httpAny, 404))

	httpAny.StatusCodeMin = pointer.ToInt(200)
	httpAny.StatusCodeMax = pointer.ToInt(399)
	assert.True(t, checkStatusCode(httpAny, 301))
	assert.False(t, checkStatusCode(httpAny, 404))
}

func TestHTTPUsecase_Match(t *testing.T) {
	httpRaw := &models.HTTPRawParams{}
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

	httpRaw.Regex = `url (\d)`
	match, substring = matchRegex(httpRaw, "test 2")
	assert.False(t, match)
	assert.Equal(t, "", substring)
}

func TestHTTPUsecase_LookupKey_Json(t *testing.T) {
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
	httpFormatted := &models.HTTPFormattedParams{}
	httpFormatted.Key = `bloc1."bloc.2".[0].value`

	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	if assert.NoError(t, err) {
		found, value := lookupKey(httpFormatted, data)
		assert.True(t, found)
		assert.Equal(t, "YEAH !!", value)
	}
}

func TestHTTPUsecase_LookupKey_Yaml(t *testing.T) {
	input := `
bloc1:
  bloc.2:
    - name: test1
      value: "YEAH !!"
    - name: test2
      value: "NOOO !!"
`
	httpYaml := &models.HTTPFormattedParams{}
	httpYaml.Key = `bloc1."bloc.2".[0].value`

	var data interface{}
	err := yaml.Unmarshal([]byte(input), &data)
	if assert.NoError(t, err) {
		found, value := lookupKey(httpYaml, data)
		assert.True(t, found)
		assert.Equal(t, "YEAH !!", value)
	}
}

func TestHTTPUsecase_LookupKey_MissingKey(t *testing.T) {
	input := `
bloc1:
  bloc.2:
    - name: test1
      value: "YEAH !!"
    - name: test2
      value: "NOOO !!"
`
	httpYaml := &models.HTTPFormattedParams{}
	httpYaml.Key = `bloc1."bloc.2".[0].value2`

	var data interface{}
	err := yaml.Unmarshal([]byte(input), &data)
	if assert.NoError(t, err) {
		found, _ := lookupKey(httpYaml, data)
		assert.False(t, found)
	}
}
