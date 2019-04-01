package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mErrors "github.com/jsdidierlaurent/monitoror/models/errors"

	"github.com/jsdidierlaurent/monitoror/models/tiles"
	. "github.com/jsdidierlaurent/monitoror/monitorable/port"
	"github.com/jsdidierlaurent/monitoror/monitorable/port/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/info", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestDelivery_GetPort_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	hostname := "test.com"
	port := "1234"
	ctx.QueryParams().Set("hostname", hostname)
	ctx.QueryParams().Set("port", port)

	tile := tiles.NewHealthTile(PortTileSubType)
	tile.Label = fmt.Sprintf("%s:%s", hostname, port)
	tile.Status = tiles.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Port", Anything).Return(tile, nil)
	handler := NewHttpPortHandler(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetPort(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Port", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetPort_QueryParamsError_MissingHostname(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	handler := NewHttpPortHandler(mockUsecase)
	ctx.QueryParams().Set("Port", "1234")

	// Test
	err := handler.GetPort(ctx)
	assert.Error(t, err)
	assert.IsType(t, &mErrors.QueryParamsError{}, err)
}

func TestDelivery_GetPort_QueryParamsError_MissingPort(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	handler := NewHttpPortHandler(mockUsecase)
	ctx.QueryParams().Set("hostname", "test.com")

	// Test
	err := handler.GetPort(ctx)
	assert.Error(t, err)
	assert.IsType(t, &mErrors.QueryParamsError{}, err)
}

func TestDelivery_GetPort_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	hostname := "test.com"
	ctx.QueryParams().Set("hostname", hostname)
	ctx.QueryParams().Set("port", "1234")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Port", Anything).Return(nil, errors.New("port error"))
	handler := NewHttpPortHandler(mockUsecase)

	// Test
	assert.Error(t, handler.GetPort(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Port", 1)
	mockUsecase.AssertExpectations(t)
}
