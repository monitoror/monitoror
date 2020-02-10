package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/monitorable/github"
	. "github.com/stretchr/testify/mock"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/github/mocks"
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

	tile := models.NewTile(github.GithubCountTileType)
	tile.Status = models.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Count", Anything).Return(tile, nil)
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

	tile := models.NewTile(github.GithubChecksTileType)
	tile.Status = models.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Checks", Anything).Return(tile, nil)
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
		assert.IsType(t, &models.MonitororError{}, err)
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
