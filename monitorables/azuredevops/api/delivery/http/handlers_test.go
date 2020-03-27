package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/azuredevops/xxx", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	ctx.QueryParams().Set("project", "test")
	ctx.QueryParams().Set("definition", "1")

	return
}

func TestDelivery_BuildHandler_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	tile := models.NewTile(api.AzureDevOpsBuildTileType)
	tile.Status = models.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Build", Anything).Return(tile, nil)
	handler := NewAzureDevOpsDelivery(mockUsecase)

	// Expected
	j, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetBuild(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Build", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_BuildHandler_QueryParamsError_MissingGroup(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del("project")

	mockUsecase := new(mocks.Usecase)
	handler := NewAzureDevOpsDelivery(mockUsecase)

	// Test
	err := handler.GetBuild(ctx)
	assert.Error(t, err)
	assert.IsType(t, &models.MonitororError{}, err)
}

func TestDelivery_BuildHandler_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Build", Anything).Return(nil, errors.New("build error"))
	handler := NewAzureDevOpsDelivery(mockUsecase)

	// Test
	assert.Error(t, handler.GetBuild(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Build", 1)
	mockUsecase.AssertExpectations(t)
}

func TestDelivery_GetRelease_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	tile := models.NewTile(api.AzureDevOpsBuildTileType)
	tile.Status = models.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Release", Anything).Return(tile, nil)
	handler := NewAzureDevOpsDelivery(mockUsecase)

	// Expected
	j, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetRelease(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Release", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetRelease_QueryParamsError_MissingGroup(t *testing.T) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del("project")

	mockUsecase := new(mocks.Usecase)
	handler := NewAzureDevOpsDelivery(mockUsecase)

	// Test
	err := handler.GetRelease(ctx)
	assert.Error(t, err)
	assert.IsType(t, &models.MonitororError{}, err)
}

func TestDelivery_GetRelease_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Release", Anything).Return(nil, errors.New("build error"))
	handler := NewAzureDevOpsDelivery(mockUsecase)

	// Test
	assert.Error(t, handler.GetRelease(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Release", 1)
	mockUsecase.AssertExpectations(t)
}
