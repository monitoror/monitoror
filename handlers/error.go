package handlers

import (
	"net/http"

	"github.com/monitoror/monitoror/middlewares"
	"github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/models/tiles"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/labstack/echo/v4"
)

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HttpErrorHandler(err error, c echo.Context) {
	switch e := err.(type) {
	case *errors.SystemError:
		systemError(c, e)
	case *errors.TimeoutError:
		timeoutError(c, e)
	case *errors.QueryParamsError:
		queryParamsError(c, e)
	default:
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == 404 {
				// 404
				_ = c.JSON(he.Code, ApiError{
					Status:  he.Code,
					Message: "Not Found",
				})
				return
			}
		}

		_ = c.JSON(http.StatusInternalServerError, ApiError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
}

func systemError(c echo.Context, e error) {
	c.Logger().Error(e.Error())
	tile := tiles.NewErrorTile("System Error", e.Error())
	_ = c.JSON(http.StatusInternalServerError, tile)
}

func queryParamsError(c echo.Context, qpe *errors.QueryParamsError) {
	c.Logger().Error(qpe.Error())
	tile := tiles.NewErrorTile("Wrong Configuration", qpe.Error())
	_ = c.JSON(http.StatusBadRequest, tile)
}

// timeoutError return cached value from downstreamStore if exist
func timeoutError(c echo.Context, te *errors.TimeoutError) {
	sendTimeout := func(c echo.Context) {
		tile := te.Tile
		tile.Status = tiles.TimeoutStatus
		tile.Message = te.Error()
		_ = c.JSON(http.StatusOK, tile)
	}

	// Looking for TimeoutCache in echo.context
	value := c.Get(middlewares.DownstreamStoreContextKey)
	if value == nil {
		c.Logger().Warn("unable to find DownstreamStore in echo.context")
		sendTimeout(c)
		return
	}
	store, ok := value.(cache.Store)
	if !ok {
		c.Logger().Warn("unable to cast value in cache.Store")
		sendTimeout(c)
		return
	}

	//Looking for Data in DownstreamStore
	var cachedResponse cache.ResponseCache
	if err := store.Get(cache.GetKey(middlewares.CachePrefix, c.Request()), &cachedResponse); err != nil {
		sendTimeout(c)
	} else {
		// Cache found, return cached Data
		for k, vals := range cachedResponse.Header {
			for _, v := range vals {
				if c.Response().Header().Get(k) == "" {
					c.Response().Header().Add(k, v)
				}
			}
		}

		// Adding Header
		c.Response().Header().Add(middlewares.DownstreamCacheHeader, "true")

		c.Response().WriteHeader(cachedResponse.Status)
		_, _ = c.Response().Write(cachedResponse.Data)
	}
}
