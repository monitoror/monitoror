package http

import (
	"encoding/json"
	"errors"
	. "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http"
	"github.com/monitoror/monitoror/monitorable/http/mocks"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"

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
			assert.IsType(t, &models.MonitororError{}, err)
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
		ctx.QueryParams().Set("format", httpModels.JSONFormat)
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
		tileType     models.TileType
		mockFuncName string
		handlerFunc  handlerFunc
	}{
		{
			tileType:     http.HTTPStatusTileType,
			mockFuncName: "HTTPStatus",
			handlerFunc: func(handler *HTTPDelivery) func(ctx echo.Context) error {
				return handler.GetHTTPStatus
			},
		},
		{
			tileType:     http.HTTPRawTileType,
			mockFuncName: "HTTPRaw",
			handlerFunc: func(handler *HTTPDelivery) func(ctx echo.Context) error {
				return handler.GetHTTPRaw
			},
		},
		{
			tileType:     http.HTTPFormattedTileType,
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

		tile := models.NewTile(testcase.tileType)
		mockUsecase := new(mocks.Usecase)
		mockUsecase.On(testcase.mockFuncName, Anything).Return(tile, nil)
		handler := NewHTTPDelivery(mockUsecase)

		// Expected
		j, err := json.Marshal(tile)
		assert.NoError(t, err, "unable to marshal tile")

		// Test
		if assert.NoError(t, testcase.handlerFunc(handler)(ctx)) {
			assert.Equal(t, StatusOK, res.Code)
			assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
			mockUsecase.AssertNumberOfCalls(t, testcase.mockFuncName, 1)
			mockUsecase.AssertExpectations(t)
		}
	}
}
