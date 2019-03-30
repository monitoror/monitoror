//+build faker

package service

import (
	"fmt"

	"github.com/jsdidierlaurent/monitoror/cli/version"
	"github.com/jsdidierlaurent/monitoror/config"
	"github.com/jsdidierlaurent/monitoror/handlers"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/delivery/http"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/usecase"

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
	// CORS
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	// ----- Routes -----
	v1 := e.Group("/api/v1")

	// ------------- INFO ------------- //
	infoHandler := handlers.HttpInfoHandler(config)
	v1.GET("/info", infoHandler.GetInfo)

	// ------------- PING ------------- //
	pingUsecase := usecase.NewPingUsecase()
	pingHandler := http.NewHttpPingHandler(pingUsecase)
	v1.GET("/ping", pingHandler.GetPing)

	// Start service
	version.Version = "x.x.x-faker"
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Port)))
}
