package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlekSi/pointer"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api"
	"github.com/monitoror/monitoror/monitorables/github/api/mocks"
	"github.com/monitoror/monitoror/monitorables/github/api/models"

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

func TestDelivery_GetCount_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("query", "test")

	tile := coreModels.NewTile(api.GithubCountTileType)
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Count", &models.CountParams{Query: "test"}).Return(tile, nil)
	handler := NewGithubDelivery(mockUsecase)

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
	handler := NewGithubDelivery(mockUsecase)

	// Test
	err := handler.GetCount(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.IsType(t, &coreModels.MonitororError{}, err)
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
	handler := NewGithubDelivery(mockUsecase)

	// Test
	err := handler.GetCount(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		mockUsecase.AssertNumberOfCalls(t, "Count", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetChecks_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("owner", "test")
	ctx.QueryParams().Set("repository", "test")
	ctx.QueryParams().Set("ref", "master")

	tile := coreModels.NewTile(api.GithubChecksTileType)
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Checks", &models.ChecksParams{Owner: "test", Repository: "test", Ref: "master"}).Return(tile, nil)
	handler := NewGithubDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetChecks(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Checks", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetChecks_MissingParams(t *testing.T) {
	// Init
	ctx, res := initEcho()

	mockUsecase := new(mocks.Usecase)
	handler := NewGithubDelivery(mockUsecase)

	// Test
	err := handler.GetChecks(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		mockUsecase.AssertNumberOfCalls(t, "Checks", 0)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetChecks_Error(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("owner", "test")
	ctx.QueryParams().Set("repository", "test")
	ctx.QueryParams().Set("ref", "master")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Checks", Anything).Return(nil, errors.New("build error"))
	handler := NewGithubDelivery(mockUsecase)

	// Test
	err := handler.GetChecks(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		mockUsecase.AssertNumberOfCalls(t, "Checks", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetPullRequest_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("owner", "test")
	ctx.QueryParams().Set("repository", "test")
	ctx.QueryParams().Set("id", "10")

	tile := coreModels.NewTile(api.GithubPullRequestTileType)
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("PullRequest", &models.PullRequestParams{Owner: "test", Repository: "test", ID: pointer.ToInt(10)}).Return(tile, nil)
	handler := NewGithubDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetPullRequest(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "PullRequest", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetPullRequest_MissingParams(t *testing.T) {
	// Init
	ctx, res := initEcho()

	mockUsecase := new(mocks.Usecase)
	handler := NewGithubDelivery(mockUsecase)

	// Test
	err := handler.GetPullRequest(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.IsType(t, &coreModels.MonitororError{}, err)
		mockUsecase.AssertNumberOfCalls(t, "PullRequest", 0)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetPullRequest_Error(t *testing.T) {
	// Init
	ctx, res := initEcho()

	ctx.QueryParams().Set("owner", "test")
	ctx.QueryParams().Set("repository", "test")
	ctx.QueryParams().Set("id", "10")

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("PullRequest", Anything).Return(nil, errors.New("build error"))
	handler := NewGithubDelivery(mockUsecase)

	// Test
	err := handler.GetPullRequest(ctx)
	if assert.Error(t, err) {
		assert.Equal(t, http.StatusOK, res.Code)
		mockUsecase.AssertNumberOfCalls(t, "PullRequest", 1)
		mockUsecase.AssertExpectations(t)
	}
}
