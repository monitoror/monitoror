package errors

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/stretchr/testify/mock"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/jsdidierlaurent/echo-middleware/cache/mocks"
	"github.com/labstack/echo/v4"
	"github.com/monitoror/monitoror/middlewares"
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/stretchr/testify/assert"
)

func initErrorEcho() (ctx echo.Context, res *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/error", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	return
}

func TestTimeoutError_WithoutCacheStore(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()

	// Parameters
	tile := tiles.NewHealthTile("TEST").Tile
	err := NewTimeoutError(tile, "service is burning")

	// Expected
	expectedTile := tile
	expectedTile.Status = tiles.WarningStatus
	expectedTile.Message = err.Error()
	j, e := json.Marshal(expectedTile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	err.Send(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
}

func TestTimeoutError_WithCastErrorOnGetCacheStore(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()
	ctx.Set(middlewares.DownstreamStoreContextKey, "store")

	// Parameters
	tile := tiles.NewHealthTile("TEST")
	err := NewTimeoutError(tile.Tile, "service is burning")

	// Expected
	expectedTile := tile
	expectedTile.Status = tiles.WarningStatus
	expectedTile.Message = err.Error()
	j, e := json.Marshal(expectedTile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	err.Send(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
}

func TestTimeoutError_CacheMiss(t *testing.T) {
	// Init
	ctx, res := initErrorEcho()
	mockStore := new(mocks.Store)
	mockStore.On("Get", AnythingOfType("string"), Anything).Return(cache.ErrCacheMiss)
	ctx.Set(middlewares.DownstreamStoreContextKey, mockStore)

	// Parameters
	tile := tiles.NewHealthTile("TEST")
	err := NewTimeoutError(tile.Tile, "service is burning")

	// Expected
	expectedTile := tile
	expectedTile.Status = tiles.WarningStatus
	expectedTile.Message = err.Error()
	j, e := json.Marshal(expectedTile)
	assert.NoError(t, e, "unable to marshal tile")

	// Test
	err.Send(ctx)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, string(j), strings.TrimSpace(res.Body.String()))
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
	te := NewTimeoutError(tiles.NewHealthTile("TEST").Tile, "service is burning")

	// Test
	te.Send(ctx)
	header.Add("Timeout-Recover", "true")

	assert.Equal(t, status, res.Code)
	assert.Equal(t, header, res.Header())
	assert.Equal(t, body, res.Body.String())
	mockStore.AssertNumberOfCalls(t, "Get", 1)
	mockStore.AssertExpectations(t)
}
