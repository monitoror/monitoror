package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/pingdom/api"
	"github.com/monitoror/monitoror/monitorables/pingdom/api/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/pingdom/check", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	ctx.QueryParams().Set("id", "123456")

	return
}

func TestDelivery_GetCheck_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	tile := coreModels.NewTile(api.PingdomCheckTileType)
	tile.Label = "check 1"
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Check", Anything).Return(tile, nil)
	handler := NewPingdomDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetCheck(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Check", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetCheck_QueryParamsError(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del("id")

	mockUsecase := new(mocks.Usecase)
	handler := NewPingdomDelivery(mockUsecase)

	// Test
	err := handler.GetCheck(ctx)
	assert.Error(t, err)
	assert.IsType(t, &coreModels.MonitororError{}, err)
}

func TestDelivery_GetCheck_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Check", Anything).Return(nil, errors.New("boom"))
	handler := NewPingdomDelivery(mockUsecase)

	// Test
	assert.Error(t, handler.GetCheck(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Check", 1)
	mockUsecase.AssertExpectations(t)
}
