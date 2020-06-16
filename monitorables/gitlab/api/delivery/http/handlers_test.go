package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/mocks"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"

	"github.com/AlekSi/pointer"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/test", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestDelivery_GetCountIssues_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("projectId", "10")
	ctx.QueryParams().Set("query", "test")

	tile := coreModels.NewTile(api.GitlabCountIssuesTileType)
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("CountIssues", &models.IssuesParams{ProjectID: pointer.ToInt(10)}).Return(tile, nil)
	handler := NewGitlabDelivery(mockUsecase)

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

func TestDelivery_GetCountIssues_Error(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("query", "test")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("CountIssues", Anything).Return(nil, errors.New("build error"))
	handler := NewGitlabDelivery(mockUsecase)

	// Test
	err := handler.GetCountIssues(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		mockUsecase.AssertNumberOfCalls(t, "CountIssues", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetPipeline_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("projectId", "10")
	ctx.QueryParams().Set("ref", "master")

	tile := coreModels.NewTile(api.GitlabPipelineTileType)
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Pipeline", &models.PipelineParams{ProjectID: pointer.ToInt(10), Ref: "master"}).Return(tile, nil)
	handler := NewGitlabDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetPipeline(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Pipeline", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetChecks_MissingParams(t *testing.T) {
	// Init
	ctx, res := initEcho()

	mockUsecase := new(mocks.Usecase)
	handler := NewGitlabDelivery(mockUsecase)

	// Test
	err := handler.GetPipeline(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		mockUsecase.AssertNumberOfCalls(t, "Pipeline", 0)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetChecks_Error(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("projectId", "10")
	ctx.QueryParams().Set("ref", "master")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Pipeline", Anything).Return(nil, errors.New("build error"))
	handler := NewGitlabDelivery(mockUsecase)

	// Test
	err := handler.GetPipeline(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		mockUsecase.AssertNumberOfCalls(t, "Pipeline", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetMergeRequest_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("projectId", "10")
	ctx.QueryParams().Set("id", "10")

	tile := coreModels.NewTile(api.GitlabMergeRequestTileType)
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("MergeRequest", &models.MergeRequestParams{ProjectID: pointer.ToInt(10), ID: pointer.ToInt(10)}).Return(tile, nil)
	handler := NewGitlabDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetMergeRequest(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "MergeRequest", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetMergeRequest_MissingParams(t *testing.T) {
	// Init
	ctx, res := initEcho()

	mockUsecase := new(mocks.Usecase)
	handler := NewGitlabDelivery(mockUsecase)

	// Test
	err := handler.GetMergeRequest(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		mockUsecase.AssertNumberOfCalls(t, "MergeRequest", 0)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetMergeRequest_Error(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("projectId", "10")
	ctx.QueryParams().Set("id", "10")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("MergeRequest", Anything).Return(nil, errors.New("build error"))
	handler := NewGitlabDelivery(mockUsecase)

	// Test
	err := handler.GetMergeRequest(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		mockUsecase.AssertNumberOfCalls(t, "MergeRequest", 1)
		mockUsecase.AssertExpectations(t)
	}
}
