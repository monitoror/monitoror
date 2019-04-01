package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mErrors "github.com/jsdidierlaurent/monitoror/models/errors"
	"github.com/jsdidierlaurent/monitoror/models/tiles"
	. "github.com/jsdidierlaurent/monitoror/monitorable/ping"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/mocks"

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

func TestDelivery_GetPing_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	hostname := "test.com"
	ctx.QueryParams().Set("hostname", hostname)

	tile := tiles.NewHealthTile(PingTileSubType)
	tile.Label = hostname
	tile.Status = tiles.SuccessStatus
	tile.Message = "1s"

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Ping", Anything).Return(tile, nil)
	handler := NewHttpPingHandler(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetPing(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Ping", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetPing_QueryParamsError(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	handler := NewHttpPingHandler(mockUsecase)

	// Test
	err := handler.GetPing(ctx)
	assert.Error(t, err)
	assert.IsType(t, &mErrors.QueryParamsError{}, err)
}

func TestDelivery_GetPing_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	hostname := "test.com"
	ctx.QueryParams().Set("hostname", hostname)

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Ping", Anything).Return(nil, errors.New("ping error"))
	handler := NewHttpPingHandler(mockUsecase)

	// Test
	assert.Error(t, handler.GetPing(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Ping", 1)
	mockUsecase.AssertExpectations(t)
}
