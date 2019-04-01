//+build !faker

package service

import (
	"fmt"

	_pingDelivery "github.com/jsdidierlaurent/monitoror/monitorable/ping/delivery/http"
	_pingRepository "github.com/jsdidierlaurent/monitoror/monitorable/ping/repository"
	_pingUsecase "github.com/jsdidierlaurent/monitoror/monitorable/ping/usecase"
	_portDelivery "github.com/jsdidierlaurent/monitoror/monitorable/port/delivery/http"
	_portRepository "github.com/jsdidierlaurent/monitoror/monitorable/port/repository"
	_portUsecase "github.com/jsdidierlaurent/monitoror/monitorable/port/usecase"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/jsdidierlaurent/monitoror/config"
	"github.com/jsdidierlaurent/monitoror/handlers"
	"github.com/jsdidierlaurent/monitoror/middlewares"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func Start(config *config.Config) {
	e := echo.New()
	e.HideBanner = true

	//  ----- Logger -----
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("⇨ ${time_rfc3339} [${level}]")
		l.SetLevel(log.INFO)
	}

	// ----- Errors Handler -----
	e.HTTPErrorHandler = handlers.HttpErrorHandler

	//  ----- Middlewares -----
	// Recover (don't panic 😎)
	e.Use(echoMiddleware.Recover())
	// Log requests
	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: `⇨ ${time_rfc3339} [REQUEST] ${method} ${uri} status:${status} error:"${error}"` + "\n",
	}))
	// Cache
	cm := middlewares.NewCacheMiddleware(config) // Used as Handler wrapper in routes
	e.Use(cm.DownstreamStoreMiddleware())
	// CORS
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	// ----- Routes -----
	v1 := e.Group("/api/v1")

	// ------------- INFO ------------- //
	infoHandler := handlers.HttpInfoHandler(config)
	v1.GET("/info", cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoHandler.GetInfo))

	// ------------- PING ------------- //
	pingRepo := _pingRepository.NewNetworkPingRepository(config)
	pingUC := _pingUsecase.NewPingUsecase(pingRepo)
	pingHandler := _pingDelivery.NewHttpPingHandler(pingUC)
	v1.GET("/ping", cm.UpstreamCacheHandler(pingHandler.GetPing))

	// ------------- PORT ------------- //
	portRepo := _portRepository.NewNetworkPortRepository(config)
	portUC := _portUsecase.NewPortUsecase(portRepo)
	portHandler := _portDelivery.NewHttpPortHandler(portUC)
	v1.GET("/port", cm.UpstreamCacheHandler(portHandler.GetPort))

	// Start service
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Port)))
}
