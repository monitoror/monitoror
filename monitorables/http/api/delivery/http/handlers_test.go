package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlekSi/pointer"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/mocks"
	"github.com/monitoror/monitoror/monitorables/http/api/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

type handlerFunc func(handler *HTTPDelivery) func(ctx echo.Context) error

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/http", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestQueryParams_HTTPStatusParams(t *testing.T) {
	ctx, _ := initEcho()
	ctx.QueryParams().Set("url", "http://monitoror.example.com")
	ctx.QueryParams().Set("statusCodeMin", "300")
	ctx.QueryParams().Set("statusCodeMax", "400")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("HTTPStatus", &models.HTTPStatusParams{
		URL:           "http://monitoror.example.com",
		StatusCodeMin: pointer.ToInt(300),
		StatusCodeMax: pointer.ToInt(400),
	}).Return(nil, nil)
	handler := NewHTTPDelivery(mockUsecase)
	assert.NoError(t, handler.GetHTTPStatus(ctx))
}

func TestQueryParams_HTTPRawParams(t *testing.T) {
	ctx, _ := initEcho()
	ctx.QueryParams().Set("url", "http://monitoror.example.com")
	ctx.QueryParams().Set("regex", "test")
	ctx.QueryParams().Set("statusCodeMin", "300")
	ctx.QueryParams().Set("statusCodeMax", "400")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("HTTPRaw", &models.HTTPRawParams{
		URL:           "http://monitoror.example.com",
		Regex:         "test",
		StatusCodeMin: pointer.ToInt(300),
		StatusCodeMax: pointer.ToInt(400),
	}).Return(nil, nil)
	handler := NewHTTPDelivery(mockUsecase)
	assert.NoError(t, handler.GetHTTPRaw(ctx))
}

func TestQueryParams_HTTPFormattedParams(t *testing.T) {
	ctx, _ := initEcho()
	ctx.QueryParams().Set("url", "http://monitoror.example.com")
	ctx.QueryParams().Set("key", "key")
	ctx.QueryParams().Set("format", "JSON")
	ctx.QueryParams().Set("regex", "test")
	ctx.QueryParams().Set("statusCodeMin", "300")
	ctx.QueryParams().Set("statusCodeMax", "400")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("HTTPFormatted", &models.HTTPFormattedParams{
		URL:           "http://monitoror.example.com",
		Regex:         "test",
		Key:           "key",
		Format:        models.JSONFormat,
		StatusCodeMin: pointer.ToInt(300),
		StatusCodeMax: pointer.ToInt(400),
	}).Return(nil, nil)
	handler := NewHTTPDelivery(mockUsecase)
	assert.NoError(t, handler.GetHTTPFormatted(ctx))
}

func Test_httpHttpDelivery_GetHttp_MissingParams(t *testing.T) {
	// init tests cases
	testcases := []handlerFunc{
		func(handler *HTTPDelivery) func(ctx echo.Context) error {
			return handler.GetHTTPStatus
		},
		func(handler *HTTPDelivery) func(ctx echo.Context) error {
			return handler.GetHTTPRaw
		},
		func(handler *HTTPDelivery) func(ctx echo.Context) error {
			return handler.GetHTTPFormatted
		},
	}

	// tests
	for _, handlerFunc := range testcases {
		ctx, _ := initEcho()
		handler := NewHTTPDelivery(nil)

		err := handlerFunc(handler)(ctx)
		if assert.Error(t, err) {
			assert.IsType(t, &coreModels.MonitororError{}, err)
		}
	}
}

func Test_httpHttpDelivery_GetHttp_Error(t *testing.T) {
	// init tests cases
	testcases := []struct {
		mockFuncName string
		handlerFunc  handlerFunc
	}{
		{
			mockFuncName: "HTTPStatus",
			handlerFunc: func(handler *HTTPDelivery) func(ctx echo.Context) error {
				return handler.GetHTTPStatus
			},
		},
		{
			mockFuncName: "HTTPRaw",
			handlerFunc: func(handler *HTTPDelivery) func(ctx echo.Context) error {
				return handler.GetHTTPRaw
			},
		},
		{
			mockFuncName: "HTTPFormatted",
			handlerFunc: func(handler *HTTPDelivery) func(ctx echo.Context) error {
				return handler.GetHTTPFormatted
			},
		},
	}

	// tests
	for _, testcase := range testcases {
		ctx, _ := initEcho()
		ctx.QueryParams().Set("url", "http://monitoror.example.com")
		ctx.QueryParams().Set("format", models.JSONFormat)
		ctx.QueryParams().Set("regex", "(.*)")
		ctx.QueryParams().Set("key", "key")

		mockUsecase := new(mocks.Usecase)
		mockUsecase.On(testcase.mockFuncName, Anything).Return(nil, errors.New("boom"))
		handler := NewHTTPDelivery(mockUsecase)

		// Test
		if assert.Error(t, testcase.handlerFunc(handler)(ctx)) {
			mockUsecase.AssertNumberOfCalls(t, testcase.mockFuncName, 1)
			mockUsecase.AssertExpectations(t)
		}
	}
}

func Test_httpHttpDelivery_GetHttp(t *testing.T) {
	// init tests cases
	testcases := []struct {
		tileType     coreModels.TileType
		mockFuncName string
		handlerFunc  handlerFunc
	}{
		{
			tileType:     api.HTTPStatusTileType,
			mockFuncName: "HTTPStatus",
			handlerFunc: func(handler *HTTPDelivery) func(ctx echo.Context) error {
				return handler.GetHTTPStatus
			},
		},
		{
			tileType:     api.HTTPRawTileType,
			mockFuncName: "HTTPRaw",
			handlerFunc: func(handler *HTTPDelivery) func(ctx echo.Context) error {
				return handler.GetHTTPRaw
			},
		},
		{
			tileType:     api.HTTPFormattedTileType,
			mockFuncName: "HTTPFormatted",
			handlerFunc: func(handler *HTTPDelivery) func(ctx echo.Context) error {
				return handler.GetHTTPFormatted
			},
		},
	}

	// tests
	for _, testcase := range testcases {
		ctx, res := initEcho()
		ctx.QueryParams().Set("url", "http://monitoror.example.com")
		ctx.QueryParams().Set("format", "JSON")
		ctx.QueryParams().Set("regex", "(.*)")
		ctx.QueryParams().Set("key", "key")

		tile := coreModels.NewTile(testcase.tileType)
		mockUsecase := new(mocks.Usecase)
		mockUsecase.On(testcase.mockFuncName, Anything).Return(tile, nil)
		handler := NewHTTPDelivery(mockUsecase)

		// Expected
		j, err := json.Marshal(tile)
		assert.NoError(t, err, "unable to marshal tile")

		// Test
		if assert.NoError(t, testcase.handlerFunc(handler)(ctx)) {
			assert.Equal(t, http.StatusOK, res.Code)
			assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
			mockUsecase.AssertNumberOfCalls(t, testcase.mockFuncName, 1)
			mockUsecase.AssertExpectations(t)
		}
	}
}
