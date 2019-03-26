//+build debug

package service

import (
	"fmt"

	"github.com/jsdidierlaurent/monitowall/monitorable/ping/delivery/http"

	"github.com/jsdidierlaurent/monitowall/monitorable/ping/usecase"

	"github.com/jsdidierlaurent/monitowall/configs"
	"github.com/jsdidierlaurent/monitowall/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartDebug(config *configs.Config, buildInfo *configs.BuildInfo) {
	router := echo.New()

	//  ----- Middlewares -----
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

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
	v1.GET("/info", infoHandler.GetInfo)

	// ------------- PING ------------- //
	pingUsecase := usecase.NewPingUsecaseDebug()
	pingHandler := http.NewHttpPingHandler(pingUsecase)
	v1.GET("/ping", pingHandler.GetPing)

	// Start service
	router.Logger.Fatal(router.Start(fmt.Sprintf(":%d", config.Port)))
}
