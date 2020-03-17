package router

import (
	"fmt"

	"github.com/monitoror/monitoror/service/options"

	"github.com/monitoror/monitoror/middlewares"

	"github.com/labstack/echo/v4"
)

type (
	MonitorableRouter interface {
		Group(path, variant string) MonitorableGroup
	}
	MonitorableGroup interface {
		GET(path string, handlerFunc echo.HandlerFunc, options ...options.RouterOption) *echo.Route
	}

	router struct {
		apiVersion      *echo.Group
		cacheMiddleware *middlewares.CacheMiddleware
	}

	group struct {
		router *router
		group  *echo.Group
	}
)

func NewMonitorableRouter(apiVersion *echo.Group, cm *middlewares.CacheMiddleware) MonitorableRouter {
	return &router{apiVersion: apiVersion, cacheMiddleware: cm}
}

func (r *router) Group(path, variant string) MonitorableGroup {
	return &group{router: r, group: r.apiVersion.Group(fmt.Sprintf(`%s/%s`, path, variant))}
}

func (g *group) GET(path string, handlerFunc echo.HandlerFunc, opts ...options.RouterOption) *echo.Route {
	routerSettings := options.ApplyOptions(opts...)

	handler := handlerFunc
	if !routerSettings.NoCache {
		if routerSettings.CustomCacheExpiration != nil {
			handler = g.router.cacheMiddleware.UpstreamCacheHandlerWithExpiration(*routerSettings.CustomCacheExpiration, handlerFunc)
		} else {
			handler = g.router.cacheMiddleware.UpstreamCacheHandler(handlerFunc)
		}
	}

	return g.group.GET(path, handler, routerSettings.Middlewares...)
}
