package main

import (
	"github.com/jsdidierlaurent/monitowall/handler"
	"github.com/jsdidierlaurent/monitowall/route"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	router := route.Init()

	// Set Bundle MiddleWare
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	// Error Handler
	router.HTTPErrorHandler = handler.JSONErrorHandler

	router.Logger.Fatal(router.Start(":1323"))
}
