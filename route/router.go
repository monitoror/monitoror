package route

import (
	"github.com/jsdidierlaurent/monitowall/api/ping"
	"github.com/jsdidierlaurent/monitowall/model"
	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {

	e := echo.New()

	// Routes
	v1 := e.Group("/api/v1")

	// ------------- PING ------------- //
	pingHandler := ping.NewHandler(model.NewPingModel())
	v1.GET("/ping", pingHandler.GetPing)

	return e
}
