package router

import (
	"testing"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/service/middlewares"
	"github.com/monitoror/monitoror/service/options"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitorableRouter(t *testing.T) {
	// Init
	g := echo.New().Group("/api/v1")
	cacheMiddleware := middlewares.NewCacheMiddleware(cache.NewGoCacheStore(time.Minute, time.Second), time.Minute, time.Minute)
	monitorableRouter := NewMonitorableRouter(g, cacheMiddleware)
	handler := func(context echo.Context) error { return nil }

	routeGroup := monitorableRouter.Group("/test", coreModels.DefaultVariantName)
	test1 := routeGroup.GET("/test1", handler)
	test2 := routeGroup.GET("/test2", handler, options.WithNoCache())
	test3 := routeGroup.GET("/test3", handler, options.WithCustomCacheExpiration(cache.NEVER))
	test4 := routeGroup.GET("/test4", handler, options.WithMiddlewares(echoMiddleware.AddTrailingSlash()))

	assert.Equal(t, "/api/v1/test/default/test1", test1.Path)
	assert.Equal(t, "/api/v1/test/default/test2", test2.Path)
	assert.Equal(t, "/api/v1/test/default/test3", test3.Path)
	assert.Equal(t, "/api/v1/test/default/test4", test4.Path)
}
