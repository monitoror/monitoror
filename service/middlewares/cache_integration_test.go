package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/monitoror/monitoror/service/handlers"

	"github.com/monitoror/monitoror/models"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

// /!\ this is an integration test /!\
// Note : It may be necessary to separate them from unit tests

func TestCacheMiddleware(t *testing.T) {
	var timeout bool = false
	var address string = ""

	// test server
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = handlers.HTTPErrorHandler

	store := cache.NewGoCacheStore(time.Minute*5, time.Millisecond*20)
	cacheMiddleware := NewCacheMiddleware(store, time.Second, time.Millisecond*20)
	e.Use(cacheMiddleware.DownstreamStoreMiddleware())

	e.GET("/test", cacheMiddleware.UpstreamCacheHandler(func(c echo.Context) error {
		if timeout {
			return &models.MonitororError{Err: context.DeadlineExceeded, Tile: &models.Tile{}}
		}
		return c.JSON(200, `Hello world`)
	}))

	// Start server
	go func() {
		_ = e.Start(fmt.Sprintf(":%d", 0))
	}()

	// Wait until echo start
	for range time.Tick(time.Millisecond * 10) {
		if e.Listener != nil {
			address = strings.Replace(e.Listener.Addr().String(), "[::]", "http://localhost", 1)
			break
		}
	}

	// TEST
	url := fmt.Sprintf("%s/test", address)
	resp, err := http.Get(url)
	if assert.NoError(t, err) {
		assert.Equal(t, 200, resp.StatusCode)
		assert.Empty(t, resp.Header.Get("Last-Modified"))
		_ = resp.Body.Close()
	}

	resp, err = http.Get(url)
	if assert.NoError(t, err) {
		assert.Equal(t, 200, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Last-Modified"))
		_ = resp.Body.Close()
	}

	// Wait until upstream cache was clean
	time.Sleep(time.Millisecond * 25)

	timeout = true
	resp, err = http.Get(url)
	if assert.NoError(t, err) {
		assert.Equal(t, 200, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get(models.DownstreamCacheHeader))
		_ = resp.Body.Close()
	}

	// Close server
	_ = e.Close()
}
