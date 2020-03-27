package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/travisci/api"
	"github.com/monitoror/monitoror/monitorables/travisci/api/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/travisci/build", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	ctx.QueryParams().Set("owner", "test")
	ctx.QueryParams().Set("repository", "test")
	ctx.QueryParams().Set("branch", "master")

	return
}

func missingParam(t *testing.T, param string) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del(param)

	mockUsecase := new(mocks.Usecase)
	handler := NewTravisCIDelivery(mockUsecase)

	// Test
	err := handler.GetBuild(ctx)
	assert.Error(t, err)
	assert.IsType(t, &coreModels.MonitororError{}, err)
}

func TestDelivery_GetBuild_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	tile := coreModels.NewTile(api.TravisCIBuildTileType)
	tile.Status = coreModels.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Build", Anything).Return(tile, nil)
	handler := NewTravisCIDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.GetBuild(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Build", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_GetBuild_QueryParamsError_MissingGroup(t *testing.T) {
	missingParam(t, "owner")
}

func TestDelivery_GetBuild_QueryParamsError_MissingRepository(t *testing.T) {
	missingParam(t, "repository")
}

func TestDelivery_GetBuild_QueryParamsError_MissingBranch(t *testing.T) {
	missingParam(t, "branch")
}

func TestDelivery_GetBuild_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Build", Anything).Return(nil, errors.New("ping error"))
	handler := NewTravisCIDelivery(mockUsecase)

	// Test
	assert.Error(t, handler.GetBuild(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Build", 1)
	mockUsecase.AssertExpectations(t)
}
