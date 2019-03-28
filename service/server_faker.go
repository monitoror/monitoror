//+build faker

package service

import (
	"fmt"

	"github.com/jsdidierlaurent/monitoror/cli/version"
	"github.com/jsdidierlaurent/monitoror/config"
	"github.com/jsdidierlaurent/monitoror/handlers"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/delivery/http"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/usecase"

	"github.com/labstack/echo"
	"github.com/labstack/echo/v4/middleware"
)

func Start(config *config.Config) {
	e := echo.New()
	e.HideBanner = true

	//  ----- Middlewares -----
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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
	v1.GET("/info", infoHandler.GetInfo)

	// ------------- PING ------------- //
	pingUsecase := usecase.NewPingUsecase()
	pingHandler := http.NewHttpPingHandler(pingUsecase)
	v1.GET("/ping", pingHandler.GetPing)

	// Start service
	version.Version = "x.x.x-faker"
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Port)))
}
