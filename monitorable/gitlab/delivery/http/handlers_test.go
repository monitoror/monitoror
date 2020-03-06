package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/monitorable/gitlab"
	. "github.com/stretchr/testify/mock"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/gitlab/mocks"
	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/test", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestDelivery_GetCount_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("query", "test")

	tile := models.NewTile(gitlab.GitlabCountTileType)
	tile.Status = models.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Count", Anything).Return(tile, nil)
	handler := NewGitlabDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetCount(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Count", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetCount_MissingParams(t *testing.T) {
	// Init
	ctx, res := initEcho()

	mockUsecase := new(mocks.Usecase)
	handler := NewGitlabDelivery(mockUsecase)

	// Test
	err := handler.GetCount(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.IsType(t, &models.MonitororError{}, err)
		mockUsecase.AssertNumberOfCalls(t, "Count", 0)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetCount_Error(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("query", "test")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Count", Anything).Return(nil, errors.New("build error"))
	handler := NewGitlabDelivery(mockUsecase)

	// Test
	err := handler.GetCount(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		mockUsecase.AssertNumberOfCalls(t, "Count", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetPipelines_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("repository", "test")
	ctx.QueryParams().Set("ref", "master")

	tile := models.NewTile(gitlab.GitlabPipelinesTileType)
	tile.Status = models.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Pipelines", Anything).Return(tile, nil)
	handler := NewGitlabDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetPipelines(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Pipelines", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetPipelines_MissingParams(t *testing.T) {
	// Init
	ctx, res := initEcho()

	mockUsecase := new(mocks.Usecase)
	handler := NewGitlabDelivery(mockUsecase)

	// Test
	err := handler.GetPipelines(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.IsType(t, &models.MonitororError{}, err)
		mockUsecase.AssertNumberOfCalls(t, "Pipelines", 0)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetPipelines_Error(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("repository", "test")
	ctx.QueryParams().Set("ref", "master")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Pipelines", Anything).Return(nil, errors.New("build error"))
	handler := NewGitlabDelivery(mockUsecase)

	// Test
	err := handler.GetPipelines(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		mockUsecase.AssertNumberOfCalls(t, "Pipelines", 1)
		mockUsecase.AssertExpectations(t)
	}
}
