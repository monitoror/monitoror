package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api/mocks"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/oneuptime/monitor", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	ctx.QueryParams().Set("id", "8176df0f-5059-402b-a3d2-f7821257d724")

	return
}

func TestDelivery_GetMonitor_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	tile := coreModels.NewTile(api.OneuptimeMonitorTileType)
	tile.Label = "monitor 1"
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Monitor", &models.MonitorParams{ID: "8176df0f-5059-402b-a3d2-f7821257d724"}).Return(tile, nil)
	handler := NewOneuptimeDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetMonitor(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Monitor", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetMonitor_QueryParamsError(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del("id")

	mockUsecase := new(mocks.Usecase)
	handler := NewOneuptimeDelivery(mockUsecase)

	// Test
	err := handler.GetMonitor(ctx)
	assert.Error(t, err)
	assert.IsType(t, &coreModels.MonitororError{}, err)
}

func TestDelivery_GetMonitor_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Monitor", Anything).Return(nil, errors.New("boom"))
	handler := NewOneuptimeDelivery(mockUsecase)

	// Test
	assert.Error(t, handler.GetMonitor(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Monitor", 1)
	mockUsecase.AssertExpectations(t)
}
