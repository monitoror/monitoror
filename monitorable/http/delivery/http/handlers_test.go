package http

import (
	"encoding/json"
	"errors"
	. "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/http"
	"github.com/monitoror/monitoror/monitorable/http/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

type handlerFunc func(handler *httpHttpDelivery) func(ctx echo.Context) error

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
		func(handler *httpHttpDelivery) func(ctx echo.Context) error {
			return handler.GetHttpAny
		},
		func(handler *httpHttpDelivery) func(ctx echo.Context) error {
			return handler.GetHttpRaw
		},
		func(handler *httpHttpDelivery) func(ctx echo.Context) error {
			return handler.GetHttpJson
		},
		func(handler *httpHttpDelivery) func(ctx echo.Context) error {
			return handler.GetHttpYaml
		},
	}

	// tests
	for _, handlerFunc := range testcases {
		ctx, _ := initEcho()
		handler := NewHttpHttpDelivery(nil)

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
			mockFuncName: "HttpAny",
			handlerFunc: func(handler *httpHttpDelivery) func(ctx echo.Context) error {
				return handler.GetHttpAny
			},
		},
		{
			mockFuncName: "HttpRaw",
			handlerFunc: func(handler *httpHttpDelivery) func(ctx echo.Context) error {
				return handler.GetHttpRaw
			},
		},
		{
			mockFuncName: "HttpJson",
			handlerFunc: func(handler *httpHttpDelivery) func(ctx echo.Context) error {
				return handler.GetHttpJson
			},
		},
		{
			mockFuncName: "HttpYaml",
			handlerFunc: func(handler *httpHttpDelivery) func(ctx echo.Context) error {
				return handler.GetHttpYaml
			},
		},
	}

	// tests
	for _, testcase := range testcases {
		ctx, _ := initEcho()
		ctx.QueryParams().Set("url", "http://monitoror.test")
		ctx.QueryParams().Set("regex", "(.*)")
		ctx.QueryParams().Set("key", ".key")

		mockUsecase := new(mocks.Usecase)
		mockUsecase.On(testcase.mockFuncName, Anything).Return(nil, errors.New("boom"))
		handler := NewHttpHttpDelivery(mockUsecase)

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
		tileType     tiles.TileType
		mockFuncName string
		handlerFunc  handlerFunc
	}{
		{
			tileType:     http.HttpAnyTileType,
			mockFuncName: "HttpAny",
			handlerFunc: func(handler *httpHttpDelivery) func(ctx echo.Context) error {
				return handler.GetHttpAny
			},
		},
		{
			tileType:     http.HttpRawTileType,
			mockFuncName: "HttpRaw",
			handlerFunc: func(handler *httpHttpDelivery) func(ctx echo.Context) error {
				return handler.GetHttpRaw
			},
		},
		{
			tileType:     http.HttpJsonTileType,
			mockFuncName: "HttpJson",
			handlerFunc: func(handler *httpHttpDelivery) func(ctx echo.Context) error {
				return handler.GetHttpJson
			},
		},
		{
			tileType:     http.HttpYamlTileType,
			mockFuncName: "HttpYaml",
			handlerFunc: func(handler *httpHttpDelivery) func(ctx echo.Context) error {
				return handler.GetHttpYaml
			},
		},
	}

	// tests
	for _, testcase := range testcases {
		ctx, res := initEcho()
		ctx.QueryParams().Set("url", "http://monitoror.test")
		ctx.QueryParams().Set("regex", "(.*)")
		ctx.QueryParams().Set("key", ".key")

		tile := tiles.NewHealthTile(testcase.tileType)
		mockUsecase := new(mocks.Usecase)
		mockUsecase.On(testcase.mockFuncName, Anything).Return(tile, nil)
		handler := NewHttpHttpDelivery(mockUsecase)

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
