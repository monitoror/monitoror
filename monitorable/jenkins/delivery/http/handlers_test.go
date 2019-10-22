package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/models"

	"github.com/monitoror/monitoror/monitorable/jenkins/mocks"
	"github.com/monitoror/monitoror/monitorable/travisci"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/jenkins/build", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	ctx.QueryParams().Set("job", "master")
	ctx.QueryParams().Set("parent", "test")

	return
}

func TestDelivery_GetBuild_Success(t *testing.T) {
	// Init
	ctx, res := initEcho()

	tile := models.NewTile(travisci.TravisCIBuildTileType)
	tile.Status = models.SuccessStatus

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Build", Anything).Return(tile, nil)
	handler := NewHttpJenkinsDelivery(mockUsecase)

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
	// Init
	ctx, _ := initEcho()
	ctx.QueryParams().Del("job")

	mockUsecase := new(mocks.Usecase)
	handler := NewHttpJenkinsDelivery(mockUsecase)

	// Test
	err := handler.GetBuild(ctx)
	assert.Error(t, err)
	assert.IsType(t, &models.MonitororError{}, err)
}

func TestDelivery_GetBuild_Error(t *testing.T) {
	// Init
	ctx, _ := initEcho()

	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("Build", Anything).Return(nil, errors.New("build error"))
	handler := NewHttpJenkinsDelivery(mockUsecase)

	// Test
	assert.Error(t, handler.GetBuild(ctx))
	mockUsecase.AssertNumberOfCalls(t, "Build", 1)
	mockUsecase.AssertExpectations(t)
}
