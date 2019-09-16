package usecase

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/monitoror/monitoror/models/tiles"

	. "github.com/stretchr/testify/mock"

	"github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/monitorable/http/mocks"
	"github.com/monitoror/monitoror/monitorable/http/models"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

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
	httpJson := &models.HttpFormattedDataParams{}
	httpJson.Key = `.bloc1."bloc.2".[0].value`

	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	if assert.NoError(t, err) {
		found, value := lookupKey(httpJson, data)
		assert.True(t, found)
		assert.Equal(t, "YEAH !!", value)
	}
}

func TestHttpUsecase_LookupKey_Json_Array(t *testing.T) {
	input := `
[
	{ "value": "YEAH !!" },
	{ "value": "NOOO !!" },
	{ "value": "NOOO !!" }
]
`
	httpJson := &models.HttpFormattedDataParams{}
	httpJson.Key = `.[0].value`

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
	httpYaml := &models.HttpFormattedDataParams{}
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
	httpYaml := &models.HttpFormattedDataParams{}
	httpYaml.Key = `.bloc1."bloc.2".[0].value2`

	var data interface{}
	err := yaml.Unmarshal([]byte(input), &data)
	if assert.NoError(t, err) {
		found, _ := lookupKey(httpYaml, data)
		assert.False(t, found)
	}
}

func TestHttpAny_Success(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 200}, nil)
	tu := NewHttpUsecase(mockRepository)

	tile, err := tu.HttpAny(&models.HttpAnyParams{Url: "toto"})
	if assert.NoError(t, err) {
		assert.Equal(t, tiles.SuccessStatus, tile.Status)
		assert.Equal(t, "toto", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHttpAny_Failure(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 404}, nil)
	tu := NewHttpUsecase(mockRepository)

	tile, err := tu.HttpAny(&models.HttpAnyParams{Url: "toto"})
	if assert.NoError(t, err) {
		assert.Equal(t, tiles.FailedStatus, tile.Status)
		assert.Equal(t, "toto", tile.Label)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHttpAny_Error(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(nil, context.DeadlineExceeded)
	tu := NewHttpUsecase(mockRepository)

	tile, err := tu.HttpAny(&models.HttpAnyParams{Url: "toto"})
	if assert.Error(t, err) {
		assert.Nil(t, tile)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHttpRaw_Success(t *testing.T) {
	body := `api errors = 123`

	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 200, Body: []byte(body)}, nil)
	tu := NewHttpUsecase(mockRepository)

	tile, err := tu.HttpRaw(&models.HttpRawParams{Url: "toto", Regex: `api errors = (\d*)`})
	if assert.NoError(t, err) {
		assert.Equal(t, tiles.SuccessStatus, tile.Status)
		assert.Equal(t, "toto", tile.Label)
		assert.Equal(t, "123", tile.Message)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHttpRaw_Failure_PatternNotFound(t *testing.T) {
	body := `api errors = 123`

	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 200, Body: []byte(body)}, nil)
	tu := NewHttpUsecase(mockRepository)

	tile, err := tu.HttpRaw(&models.HttpRawParams{Url: "toto", Regex: `api warning = (\d*)`})
	if assert.NoError(t, err) {
		assert.Equal(t, tiles.FailedStatus, tile.Status)
		assert.Equal(t, "toto", tile.Label)
		assert.Equal(t, `pattern not found "api warning = (\d*)"`, tile.Message)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHttpJson_Success(t *testing.T) {
	body := `{"key": "value"}`

	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 200, Body: []byte(body)}, nil)
	tu := NewHttpUsecase(mockRepository)

	tile, err := tu.HttpJson(&models.HttpFormattedDataParams{Url: "toto", Key: `.key`})
	if assert.NoError(t, err) {
		assert.Equal(t, tiles.SuccessStatus, tile.Status)
		assert.Equal(t, "toto", tile.Label)
		assert.Equal(t, "value", tile.Message)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHttpJson_Failure_Unmarshal(t *testing.T) {
	body := `{"key": "value`

	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 200, Body: []byte(body)}, nil)
	tu := NewHttpUsecase(mockRepository)

	tile, err := tu.HttpJson(&models.HttpFormattedDataParams{Url: "toto", Key: `.key`})
	if assert.NoError(t, err) {
		assert.Equal(t, tiles.FailedStatus, tile.Status)
		assert.Equal(t, "toto", tile.Label)
		assert.Equal(t, "unable to unmarshal content", tile.Message)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHttpJson_Failure_MissingKey(t *testing.T) {
	body := `{"key": "value"}`

	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 200, Body: []byte(body)}, nil)
	tu := NewHttpUsecase(mockRepository)

	tile, err := tu.HttpJson(&models.HttpFormattedDataParams{Url: "toto", Key: `.missingKey`})
	if assert.NoError(t, err) {
		assert.Equal(t, tiles.FailedStatus, tile.Status)
		assert.Equal(t, "toto", tile.Label)
		assert.Equal(t, `unable to lookup for key ".missingKey"`, tile.Message)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}

func TestHttpYaml_Success(t *testing.T) {
	body := `key: value`

	mockRepository := new(mocks.Repository)
	mockRepository.On("Get", AnythingOfType("string")).Return(&models.Response{StatusCode: 200, Body: []byte(body)}, nil)
	tu := NewHttpUsecase(mockRepository)

	tile, err := tu.HttpYaml(&models.HttpFormattedDataParams{Url: "toto", Key: `.key`})
	if assert.NoError(t, err) {
		assert.Equal(t, tiles.SuccessStatus, tile.Status)
		assert.Equal(t, "toto", tile.Label)
		assert.Equal(t, "value", tile.Message)
		mockRepository.AssertNumberOfCalls(t, "Get", 1)
		mockRepository.AssertExpectations(t)
	}
}
