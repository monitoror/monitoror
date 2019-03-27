package service

import (
	"fmt"

	"github.com/jsdidierlaurent/monitowall/monitorable/ping/delivery/http"

	"github.com/jsdidierlaurent/monitowall/monitorable/ping/usecase"

	"github.com/jsdidierlaurent/monitowall/monitorable/ping/repository"

	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/jsdidierlaurent/monitowall/configs"
	"github.com/jsdidierlaurent/monitowall/handlers"
	"github.com/jsdidierlaurent/monitowall/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start(config *configs.Config, buildInfo *configs.BuildInfo) {
	router := echo.New()

	//  ----- Middlewares -----
	cm := middlewares.NewCacheMiddleware(config) // Used as Handler wrapper in routes

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(cm.DownstreamStoreMiddleware())

	//  ----- CORS -----
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	// ----- Errors Handler -----
	router.HTTPErrorHandler = handlers.HttpErrorHandler

	// ----- Routes -----
	v1 := router.Group("/api/v1")

	// ------------- INFO ------------- //
	infoHandler := handlers.HttpInfoHandler(buildInfo, config)
	v1.GET("/info", cm.UpstreamCacheHandlerWithExpiration(cache.NEVER, infoHandler.GetInfo))

	// ------------- PING ------------- //
	pingRepo := repository.NewNetworkPingRepository(config)
	pingUsecase := usecase.NewPingUsecase(pingRepo)
	pingHandler := http.NewHttpPingHandler(pingUsecase)
	v1.GET("/ping", cm.UpstreamCacheHandler(pingHandler.GetPing))

	// Start service
	router.Logger.Fatal(router.Start(fmt.Sprintf(":%d", config.Port)))
}
