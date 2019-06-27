package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mErrors "github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/travisci"
	"github.com/monitoror/monitoror/monitorable/travisci/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/travisci/build", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	ctx.QueryParams().Set("group", "test")
	ctx.QueryParams().Set("repository", "test")
	ctx.QueryParams().Set("branch", "master")

	return
}

func missingParam(t *testing.T, param string) {
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del(param)

	mockUsecase := new(mocks.Usecase)
	handler := NewHttpTravisCIDelivery(mockUsecase)

	// Test
	err := handler.MonitorBuild(ctx)
	assert.Error(t, err)
	assert.IsType(t, &mErrors.QueryParamsError{}, err)
}

func TestDelivery_MonitorBuild_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	tile := tiles.NewBuildTile(travisci.TravisCIBuildTileType)
	tile.Label = "test : #master"
	tile.Status = tiles.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Build", Anything).Return(tile, nil)
	handler := NewHttpTravisCIDelivery(mockUsecase)

	// Expected
	json, err := json.Marshal(tile)
	assert.NoError(t, err, "unable to marshal tile")

	// Test
	if assert.NoError(t, handler.MonitorBuild(ctx)) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
		mockUsecase.AssertNumberOfCalls(t, "Build", 1)
		mockUsecase.AssertExpectations(t)
	}
}

func TestDelivery_MonitorBuild_QueryParamsError_MissingGroup(t *testing.T) {
	missingParam(t, "group")
}

func TestDelivery_MonitorBuild_QueryParamsError_MissingRepository(t *testing.T) {
	missingParam(t, "repository")
}

func TestDelivery_MonitorBuild_QueryParamsError_MissingBranch(t *testing.T) {
	missingParam(t, "branch")
}

func TestDelivery_GetPing_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Build", Anything).Return(nil, errors.New("ping error"))
	handler := NewHttpTravisCIDelivery(mockUsecase)

	// Test
	assert.Error(t, handler.MonitorBuild(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Build", 1)
	mockUsecase.AssertExpectations(t)
}
