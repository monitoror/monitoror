package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/youtrack/api"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/mocks"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/youtrack/count-issues", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	ctx.QueryParams().Set("query", "test")
	ctx.QueryParams().Set("countThreshold", "map[WARNING:5]")

	return
}

func TestDelivery_GetCountIssues_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	tile := coreModels.NewTile(api.YoutrackCountIssuesTileType)
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("CountIssues", &models.IssuesCountParams{Query: "test", CountThreshold: map[coreModels.TileStatus]int{coreModels.WarningStatus: 5}}).Return(tile, nil)
	handler := NewYoutrackDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetCountIssues(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "CountIssues", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetCountIssues_QueryParamsError_MissingQuery(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del("query")

	mockUsecase := new(mocks.Usecase)
	handler := NewYoutrackDelivery(mockUsecase)

	// Test
	err := handler.GetCountIssues(ctx)
	assert.Error(t, err)
	assert.IsType(t, &coreModels.MonitororError{}, err)
}

func TestDelivery_GetCountIssues_QueryParamsError_InvalidParams(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Set("PriorityFieldThreshold", "map[WARNINGS:5]")

	mockUsecase := new(mocks.Usecase)
	handler := NewYoutrackDelivery(mockUsecase)

	// Test
	err := handler.GetCountIssues(ctx)
	assert.Error(t, err)
	assert.IsType(t, &coreModels.MonitororError{}, err)
}

func TestDelivery_GetCountIssues_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("CountIssues", mock.Anything).Return(nil, errors.New("boom"))
	handler := NewYoutrackDelivery(mockUsecase)

	// Test
	assert.Error(t, handler.GetCountIssues(ctx))
	mockUsecase.AssertNumberOfCalls(t, "CountIssues", 1)
	mockUsecase.AssertExpectations(t)
}
