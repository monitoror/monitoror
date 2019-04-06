package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monitoror/monitoror/middlewares"
	mErrors "github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/models/tiles"

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
		Status:  http.StatusNotFound,
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
		Status:  http.StatusInternalServerError,
		Message: err.Error(),
	}
	json, e := json.Marshal(apiError)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestSystemError(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	message := "system error"
	err := mErrors.NewSystemError(message, nil)

	// Expected
	tile := tiles.NewErrorTile("System Error", message)
	json, e := json.Marshal(tile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestQueryParamsError(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	err := mErrors.NewQueryParamsError(nil)

	// Expected
	tile := tiles.NewErrorTile("Wrong Configuration", err.Error())
	json, e := json.Marshal(tile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestTimeoutError_WithoutCacheStore(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	err := mErrors.NewTimeoutError(tiles.NewHealthTile("TEST").Tile, "service is burning")

	// Expected
	tile := err.Tile
	tile.Status = tiles.TimeoutStatus
	tile.Message = err.Error()
	json, e := json.Marshal(tile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestTimeoutError_WithCastErrorOnGetCacheStore(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()
	ctx.Set(middlewares.DownstreamStoreContextKey, "store")

	// Parameters
	err := mErrors.NewTimeoutError(tiles.NewHealthTile("TEST").Tile, "service is burning")

	// Expected
	tile := err.Tile
	tile.Status = tiles.TimeoutStatus
	tile.Message = err.Error()
	json, e := json.Marshal(tile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
}

func TestTimeoutError_CacheMiss(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()
	mockStore := new(mocks.Store)
	mockStore.On("Get", AnythingOfType("string"), Anything).Return(cache.ErrCacheMiss)
	ctx.Set(middlewares.DownstreamStoreContextKey, mockStore)

	// Parameters
	err := mErrors.NewTimeoutError(tiles.NewHealthTile("TEST").Tile, "service is burning")

	// Expected
	tile := err.Tile
	tile.Status = tiles.TimeoutStatus
	tile.Message = err.Error()
	json, e := json.Marshal(tile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	HttpErrorHandler(err, ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(json), strings.TrimSpace(res.Body.String()))
	mockStore.AssertNumberOfCalls(t, "Get", 1)
	mockStore.AssertExpectations(t)
}

func TestTimeoutError_Success(t *testing.T) {
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
	te := mErrors.NewTimeoutError(tiles.NewHealthTile("TEST").Tile, "service is burning")

	// Test
	HttpErrorHandler(te, ctx)
	header.Add("Timeout-Recover", "true")

	assert.Equal(t, status, res.Code)
	assert.Equal(t, header, res.Header())
	assert.Equal(t, body, res.Body.String())
	mockStore.AssertNumberOfCalls(t, "Get", 1)
	mockStore.AssertExpectations(t)
}
