package handlers

import (
	"net/http"

	"github.com/monitoror/monitoror/models"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/common/log"
)

type (
	APIError struct {
		Code    int    `json:"status"`
		Message string `json:"message"`
	}
)

func HTTPErrorHandler(err error, ctx echo.Context) {
	switch e := err.(type) {
	case *models.MonitororError:
		err = handleMonitororError(e, ctx)
	default:
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == http.StatusNotFound {
				// 404
				_ = ctx.JSON(he.Code, APIError{
					Code:    he.Code,
					Message: "Not Found",
				})
				return
			}
		}
	}

	if err != nil {
		_ = ctx.JSON(http.StatusInternalServerError, APIError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
}

func handleMonitororError(me *models.MonitororError, ctx echo.Context) error {
	// No tile set, forward error
	if me.Tile == nil {
		return me
	}

	// Check if error was timeout and check cache
	if me.Timeout() {
		// If cache found, reply cache and exit
		if found := cacheMiddleware(ctx); found {
			return nil
		}

		// Cache not found, reply Timeout based on Tile
		me.Tile.Status = models.WarningStatus
		me.Tile.Message = "timeout/host unreachable"

		_ = ctx.JSON(http.StatusOK, me.Tile)
		return nil
	}

	tile := me.Tile
	tile.Message = me.Error()
	tile.Status = me.ErrorStatus
	if tile.Status == "" {
		tile.Status = models.FailedStatus
	}

	_ = ctx.JSON(http.StatusOK, tile)
	return nil
}

// cacheMiddleware look into downstream cache and return cached value to client
func cacheMiddleware(ctx echo.Context) bool {
	// Looking for TimeoutCache in echo.context
	value := ctx.Get(models.DownstreamStoreContextKey)
	if value == nil {
		log.Errorf("unable to find DownstreamStore in echo.context")
		return false
	}

	store, ok := value.(cache.Store)
	if !ok {
		log.Errorf("unable to cast value in cache.Store")
		return false
	}

	// Looking for Data in DownstreamStore
	var cachedResponse cache.ResponseCache
	if err := store.Get(cache.GetKey(models.DownstreamStoreKeyPrefix, ctx.Request()), &cachedResponse); err != nil {
		// Cache not found, return
		return false
	}

	for k, vals := range cachedResponse.Header {
		for _, v := range vals {
			if ctx.Response().Header().Get(k) == "" {
				ctx.Response().Header().Add(k, v)
			}
		}
	}

	// Adding Header (used by cache middleware to skip "caching" this response again.
	ctx.Response().Header().Add(models.DownstreamCacheHeader, "true")

	ctx.Response().WriteHeader(cachedResponse.Status)
	_, _ = ctx.Response().Write(cachedResponse.Data)

	return true
}
