package errors

import (
	"fmt"
	"net/http"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/monitoror/monitoror/middlewares"
	"github.com/monitoror/monitoror/models/tiles"
	. "github.com/monitoror/monitoror/models/tiles"
)

type TimeoutError struct {
	Tile *Tile
}

func NewTimeoutError(tile *Tile) *TimeoutError {
	return &TimeoutError{tile}
}

// timeoutError return cached value from downstreamStore if exist
func (te *TimeoutError) Send(ctx echo.Context) {
	send := func() {
		tile := te.Tile
		tile.Status = tiles.WarningStatus
		tile.Message = te.Error()
		_ = ctx.JSON(http.StatusOK, tile)
	}

	// Looking for TimeoutCache in echo.context
	value := ctx.Get(middlewares.DownstreamStoreContextKey)
	if value == nil {
		log.Warn("unable to find DownstreamStore in echo.context")
		send()
		return
	}
	store, ok := value.(cache.Store)
	if !ok {
		log.Warn("unable to cast value in cache.Store")
		send()
		return
	}

	//Looking for Data in DownstreamStore
	var cachedResponse cache.ResponseCache
	if err := store.Get(cache.GetKey(middlewares.CachePrefix, ctx.Request()), &cachedResponse); err != nil {
		send()
	} else {
		// Cache found, return cached Data
		for k, vals := range cachedResponse.Header {
			for _, v := range vals {
				if ctx.Response().Header().Get(k) == "" {
					ctx.Response().Header().Add(k, v)
				}
			}
		}

		// Adding Header
		ctx.Response().Header().Add(middlewares.DownstreamCacheHeader, "true")

		ctx.Response().WriteHeader(cachedResponse.Status)
		_, _ = ctx.Response().Write(cachedResponse.Data)
	}
}

func (te *TimeoutError) Error() string {
	return fmt.Sprintf("timeout/host unreachable")
}
