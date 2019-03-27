package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jsdidierlaurent/monitowall/models/tiles"
	"github.com/jsdidierlaurent/monitowall/monitorable/ping/usecase"

	. "github.com/stretchr/testify/mock"

	"github.com/jsdidierlaurent/monitowall/monitorable/ping/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/info", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestDelivery_GetPing_Success(t *testing.T) {
	hostname := "test.com"

	// Init
	ctx, res := initEcho()
	ctx.QueryParams().Set("hostname", hostname)

	tile := tiles.NewHealthTile(usecase.PingTileSubType)
	tile.Label = hostname
	tile.Status = tiles.SuccessStatus
	tile.Message = "1s"

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Ping", Anything).Return(tile, nil)
	handler := NewHttpPingHandler(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)

	// Test
	assert.NoError(t, err)
	assert.NoError(t, handler.GetPing(ctx))
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestDelivery_GetPing_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Ping", Anything).Return(nil, errors.New("ping error"))
	handler := NewHttpPingHandler(mockUsecase)

	// Test
	assert.Error(t, handler.GetPing(ctx))
}
