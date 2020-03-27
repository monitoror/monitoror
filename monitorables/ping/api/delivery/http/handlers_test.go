package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	"github.com/monitoror/monitoror/monitorables/ping/api/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/info", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	ctx.QueryParams().Set("hostname", "monitoror.example.com")

	return
}

func TestDelivery_PingHandler_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	tile := coreModels.NewTile(api.PingTileType)
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Ping", Anything).Return(tile, nil)
	handler := NewPingDelivery(mockUsecase)

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

func TestDelivery_PingHandler_QueryParamsError(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del("hostname")

	mockUsecase := new(mocks.Usecase)
	handler := NewPingDelivery(mockUsecase)

	// Test
	err := handler.GetPing(ctx)
	assert.Error(t, err)
	assert.IsType(t, &coreModels.MonitororError{}, err)
}

func TestDelivery_PingHandler_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Ping", Anything).Return(nil, errors.New("ping error"))
	handler := NewPingDelivery(mockUsecase)

	// Test
	assert.Error(t, handler.GetPing(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Ping", 1)
	mockUsecase.AssertExpectations(t)
}
