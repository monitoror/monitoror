package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/jsdidierlaurent/echo-middleware/cache"
	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http"
	"github.com/monitoror/monitoror/monitorable/http/mocks"
	"github.com/monitoror/monitoror/monitorable/http/models"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v2"
)

func TestHTTPAny_WithError(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(nil, context.DeadlineExceeded)
	tu := NewHTTPUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

	tile, err := tu.HTTPAny(&models.HTTPAnyParams{URL: "toto"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHtmlAll_WithoutErrors(t *testing.T) {
	for _, testcase := range []struct {
		body            string
		usecaseFunc     func(usecase http.Usecase) (*Tile, error)
		expectedStatus  TileStatus
		expectedLabel   string
		expectedMessage string
		expectedValues  []float64
	}{
		{
			// HTTP Any
			usecaseFunc: func(usecase http.Usecase) (*Tile, error) {
				return usecase.HTTPAny(&models.HTTPAnyParams{URL: "toto"})
			},
			expectedStatus: SuccessStatus, expectedLabel: "toto",
		},
		{
			// HTTP Any with wrong status
			usecaseFunc: func(usecase http.Usecase) (*Tile, error) {
				return usecase.HTTPAny(&models.HTTPAnyParams{URL: "toto", StatusCodeMin: pointer.ToInt(400), StatusCodeMax: pointer.ToInt(499)})
			},
			expectedStatus: FailedStatus, expectedLabel: "toto", expectedMessage: "status code 200",
		},
		{
			// HTTP Raw with matched regex
			body: "errors: 28",
			usecaseFunc: func(usecase http.Usecase) (*Tile, error) {
				return usecase.HTTPRaw(&models.HTTPRawParams{URL: "toto", Regex: `errors: (\d*)`})
			},
			expectedStatus: SuccessStatus, expectedLabel: "toto", expectedValues: []float64{28},
		},
		{
			// HTTP Raw without matched regex
			body: "api call: 20",
			usecaseFunc: func(usecase http.Usecase) (*Tile, error) {
				return usecase.HTTPRaw(&models.HTTPRawParams{URL: "toto", Regex: `errors: (\d*)`})
			},
			expectedStatus: FailedStatus, expectedLabel: "toto", expectedMessage: `api call: 20`,
		},
		{
			// HTTP Json
			body: `{"key": "value"}`,
			usecaseFunc: func(usecase http.Usecase) (*Tile, error) {
				return usecase.HTTPJson(&models.HTTPJsonParams{URL: "toto", Key: ".key"})
			},
			expectedStatus: SuccessStatus, expectedLabel: "toto", expectedMessage: "value",
		},
		{
			// HTTP Json with long float
			body: `{"key": 123456789 }`,
			usecaseFunc: func(usecase http.Usecase) (*Tile, error) {
				return usecase.HTTPJson(&models.HTTPJsonParams{URL: "toto", Key: ".key"})
			},
			expectedStatus: SuccessStatus, expectedLabel: "toto", expectedValues: []float64{123456789},
		},
		{
			// HTTP Json missing key
			body: `{"key": "value"}`,
			usecaseFunc: func(usecase http.Usecase) (*Tile, error) {
				return usecase.HTTPJson(&models.HTTPJsonParams{URL: "toto", Key: ".key2"})
			},
			expectedStatus: FailedStatus, expectedLabel: "toto", expectedMessage: `unable to lookup for key ".key2"`,
		},
		{
			// HTTP Json unable to unmarshal
			body: `{"key": "value`,
			usecaseFunc: func(usecase http.Usecase) (*Tile, error) {
				return usecase.HTTPJson(&models.HTTPJsonParams{URL: "toto", Key: ".key"})
			},
			expectedStatus: FailedStatus, expectedLabel: "toto", expectedMessage: `unable to unmarshal content`,
		},
		{
			// HTTP Yaml
			body: "key: value",
			usecaseFunc: func(usecase http.Usecase) (*Tile, error) {
				return usecase.HTTPYaml(&models.HTTPYamlParams{URL: "toto", Key: ".key"})
			},
			expectedStatus: SuccessStatus, expectedLabel: "toto", expectedMessage: "value",
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
			assert.Equal(t, testcase.expectedValues, tile.Values)
			mockRepository.AssertNumberOfCalls(t, "Get", 1)
			mockRepository.AssertExpectations(t)
		}
	}
}

func TestHTTPAny_WithCache(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).
		Return(&models.Response{StatusCode: 200, Body: []byte("test with cache")}, nil)

	tu := NewHTTPUsecase(mockRepository, cache.NewGoCacheStore(time.Minute*5, time.Second), 2000)

	tile, err := tu.HTTPRaw(&models.HTTPRawParams{URL: "toto"})
	if assert.NoError(t, err) {
		assert.Equal(t, "toto", tile.Label)
		assert.Equal(t, "test with cache", tile.Message)
	}

	tile, err = tu.HTTPRaw(&models.HTTPRawParams{URL: "toto"})
	if assert.NoError(t, err) {
		assert.Equal(t, "toto", tile.Label)
		assert.Equal(t, "test with cache", tile.Message)
	}
	mockRepository.AssertNumberOfCalls(t, "Get", 1)
	mockRepository.AssertExpectations(t)
}

func TestHTTPUsecase_CheckStatusCode(t *testing.T) {
	httpAny := &models.HTTPAnyParams{}
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

	httpRaw.Regex = `toto (\d)`
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
	httpJSON := &models.HTTPJsonParams{}
	httpJSON.Key = `.bloc1."bloc.2".[0].value`

	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	if assert.NoError(t, err) {
		found, value := lookupKey(httpJSON, data)
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
	httpYaml := &models.HTTPJsonParams{}
	httpYaml.Key = `.bloc1."bloc.2".[0].value`

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
	httpYaml := &models.HTTPJsonParams{}
	httpYaml.Key = `.bloc1."bloc.2".[0].value2`

	var data interface{}
	err := yaml.Unmarshal([]byte(input), &data)
	if assert.NoError(t, err) {
		found, _ := lookupKey(httpYaml, data)
		assert.False(t, found)
	}
}
