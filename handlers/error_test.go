package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/middlewares"
	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/jenkins"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/jsdidierlaurent/echo-middleware/cache/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initErrorEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/error", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestHttpError_404(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	err := echo.NewHTTPError(http.StatusNotFound, "not found")

	// Expected
	apiError := ApiError{
		Code:    http.StatusNotFound,
		Message: "Not Found",
	}
	json, e := json.Marshal(apiError)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestHttpError_500(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	err := errors.New("boom")

	// Expected
	apiError := ApiError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
	json, e := json.Marshal(apiError)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestHttpError_MonitororError_WithoutTile(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	err := &models.MonitororError{Err: errors.New("boom"), Message: "rly big boom"}

	// Expected
	apiError := ApiError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
	json, e := json.Marshal(apiError)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestHttpError_MonitororError_WithTile(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	tile := tiles.NewBuildTile(jenkins.JenkinsBuildTileType)
	tile.Label = "test jenkins"
	err := &models.MonitororError{Err: errors.New("boom"), Tile: tile.Tile, Message: "rly big boom"}

	// Expected
	expected := tiles.NewBuildTile(jenkins.JenkinsBuildTileType)
	expected.Label = "test jenkins"
	expected.Status = tiles.FailedStatus
	expected.Message = "rly big boom"
	json, e := json.Marshal(expected)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestHttpError_MonitororError_Timeout_WithoutStore(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	tile := tiles.NewHealthTile("TEST")
	err := &models.MonitororError{Err: context.DeadlineExceeded, Tile: tile.Tile}

	// Expected
	expectedTile := tile
	expectedTile.Status = tiles.WarningStatus
	expectedTile.Message = "timeout/host unreachable"
	j, e := json.Marshal(expectedTile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
}

func TestHttpError_MonitororError_Timeout_WithWrongStore(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()
	ctx.Set(middlewares.DownstreamStoreContextKey, "store")

	// Parameters
	tile := tiles.NewHealthTile("TEST")
	err := &models.MonitororError{Err: context.DeadlineExceeded, Tile: tile.Tile}

	// Expected
	expectedTile := tile
	expectedTile.Status = tiles.WarningStatus
	expectedTile.Message = "timeout/host unreachable"
	j, e := json.Marshal(expectedTile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
}

func TestHttpError_MonitororError_Timeout_CacheMiss(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()
	mockStore := new(mocks.Store)
	mockStore.On("Get", AnythingOfType("string"), Anything).Return(cache.ErrCacheMiss)
	ctx.Set(middlewares.DownstreamStoreContextKey, mockStore)

	// Parameters
	tile := tiles.NewHealthTile("TEST")
	err := &models.MonitororError{Err: context.DeadlineExceeded, Tile: tile.Tile}

	// Expected
	expectedTile := tile
	expectedTile.Status = tiles.WarningStatus
	expectedTile.Message = "timeout/host unreachable"
	j, e := json.Marshal(expectedTile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
	mockStore.AssertNumberOfCalls(t, "Get", 1)
	mockStore.AssertExpectations(t)
}

func TestHttpError_MonitororError_Timeout_Success(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	status := http.StatusOK
	header := ctx.Request().Header
	header.Add("header", "true")
	body := "body"

	mockStore := new(mocks.Store)
	mockStore.
		On("Get", AnythingOfType("string"), AnythingOfType("*cache.ResponseCache")).
		Return(nil).
		Run(func(args Arguments) {
			arg := args.Get(1).(*cache.ResponseCache)
			arg.Data = []byte(body)
			arg.Header = header
			arg.Status = status
		})
	ctx.Set(middlewares.DownstreamStoreContextKey, mockStore)

	// Parameters
	tile := tiles.NewHealthTile("TEST")
	err := &models.MonitororError{Err: context.DeadlineExceeded, Tile: tile.Tile}

	// Test
	HttpErrorHandler(err, ctx)
	header.Add("Timeout-Recover", "true")

	assert.Equal(t, status, res.Code)
	assert.Equal(t, header, res.Header())
	assert.Equal(t, body, res.Body.String())
	mockStore.AssertNumberOfCalls(t, "Get", 1)
	mockStore.AssertExpectations(t)
}
