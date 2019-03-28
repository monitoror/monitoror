//+build !faker

package service

import (
	"fmt"

	"github.com/jsdidierlaurent/monitowall/monitorable/ping/delivery/http"

	"github.com/jsdidierlaurent/monitowall/monitorable/ping/usecase"

	"github.com/jsdidierlaurent/monitowall/monitorable/ping/repository"

	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/jsdidierlaurent/monitowall/config"
	"github.com/jsdidierlaurent/monitowall/handlers"
	"github.com/jsdidierlaurent/monitowall/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start(config *config.Config) {
	e := echo.New()
	e.HideBanner = true

	//  ----- Middlewares -----
	cm := middlewares.NewCacheMiddleware(config) // Used as Handler wrapper in routes

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(cm.DownstreamStoreMiddleware())

	//  ----- CORS -----
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	// ----- Errors Handler -----
	e.HTTPErrorHandler = handlers.HttpErrorHandler

	// ----- Routes -----
	v1 := e.Group("/api/v1")

	// ------------- INFO ------------- //
	infoHandler := handlers.HttpInfoHandler(config)
	v1.GET("/info", cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoHandler.GetInfo))

	// ------------- PING ------------- //
	pingRepo := repository.NewNetworkPingRepository(config)
	pingUsecase := usecase.NewPingUsecase(pingRepo)
	pingHandler := http.NewHttpPingHandler(pingUsecase)
	v1.GET("/ping", cm.UpstreamCacheHandler(pingHandler.GetPing))

	// Start service
	PrintBanner()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Port)))
}
