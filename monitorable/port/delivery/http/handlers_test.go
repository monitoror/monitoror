package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mErrors "github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/models/tiles"
	. "github.com/monitoror/monitoror/monitorable/port"
	"github.com/monitoror/monitoror/monitorable/port/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/info", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	ctx.QueryParams().Set("hostname", "test.com")
	ctx.QueryParams().Set("port", "1234")

	return
}

func missingParam(t *testing.T, param string) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del(param)
	mockUsecase := new(mocks.Usecase)
	handler := NewHttpPortDelivery(mockUsecase)
	// Test
	err := handler.GetPort(ctx)
	assert.Error(t, err)
	assert.IsType(t, &mErrors.QueryParamsError{}, err)
}

func TestDelivery_PortHandler_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	tile := tiles.NewHealthTile(PortTileType)
	tile.Label = "test.com:1234"
	tile.Status = tiles.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Port", Anything).Return(tile, nil)
	handler := NewHttpPortDelivery(mockUsecase)

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

func TestDelivery_PortHandler_QueryParamsError_MissingHostname(t *testing.T) {
	missingParam(t, "hostname")
}

func TestDelivery_PortHandler_QueryParamsError_MissingPort(t *testing.T) {
	missingParam(t, "port")
}

func TestDelivery_PortHandler_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Port", Anything).Return(nil, errors.New("port error"))
	handler := NewHttpPortDelivery(mockUsecase)

	// Test
	assert.Error(t, handler.GetPort(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Port", 1)
	mockUsecase.AssertExpectations(t)
}
