package handlers

import (
	"net/http"

	"github.com/jsdidierlaurent/monitowall/renderings"

	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/jsdidierlaurent/monitowall/middlewares"

	"github.com/jsdidierlaurent/monitowall/errors"
	"github.com/labstack/echo/v4"
)

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HTTPErrorHandler(err error, c echo.Context) {
	switch e := err.(type) {
	case *errors.SystemError:
		systemError(c, e)
	case *errors.TimeoutError:
		timeoutError(c, e)
	default:
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == 404 {
				// 404
				_ = c.JSON(he.Code, ApiError{
					Status:  he.Code,
					Message: "Not Found",
				})
			}
		} else {
			systemError(c, e)
		}
	}
}

func systemError(c echo.Context, e error) {
	c.Logger().Error(e.Error())
	_ = c.JSON(http.StatusInternalServerError, ApiError{
		Status:  http.StatusInternalServerError,
		Message: "System Error",
	})
}

// timeoutError return cached value from downstreamStore if exist
func timeoutError(c echo.Context, te *errors.TimeoutError) {
	// Looking for TimeoutCache in echo.context
	value := c.Get(middlewares.DownstreamStoreContextKey)
	if value == nil {
		c.Logger().Warn("unable to find DownstreamStore in echo.context")
		return
	}
	store, ok := value.(*cache.GoCacheStore)
	if !ok {
		c.Logger().Warn("unable to cast value in *cache.Store")
		return
	}

	//Looking for Data in DownstreamStore
	var cachedResponse cache.ResponseCache
	if err := store.Get(cache.GetKey(middlewares.CachePrefix, c.Request()), &cachedResponse); err != nil {
		// Missing cache, returning Timeout
		response := renderings.Response{
			Type:    te.Type,
			Status:  renderings.TimeoutStatus,
			Message: te.Message,
		}
		_ = c.JSON(http.StatusOK, response)
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
