package service

import (
	"github.com/jsdidierlaurent/monitowall/handlers"
	"github.com/jsdidierlaurent/monitowall/handlers/ping"
	"github.com/jsdidierlaurent/monitowall/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() {
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
	router.HTTPErrorHandler = handlers.JSONErrorHandler

	// ----- Routes -----
	v1 := router.Group("/api/v1")

	// ------------- PING ------------- //
	pingHandler := ping.NewHandler(models.NewPingModel())
	v1.GET("/ping", pingHandler.GetPing)


	// Start service
	router.Logger.Fatal(router.Start(":1323"))
}
