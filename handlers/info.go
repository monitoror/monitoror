package handlers

import (
	"net/http"

	"github.com/jsdidierlaurent/monitowall/config"

	"github.com/jsdidierlaurent/monitowall/middlewares"
	"github.com/jsdidierlaurent/monitowall/renderings"

	"github.com/labstack/echo/v4"
)

func GetInfo(c echo.Context) error {
	response := renderings.InfoResponse{
		BuildInfo: *c.Get(middlewares.BuildInfoContextKey).(*config.BuildInfo),
		Config:    *c.Get(middlewares.ConfigContextKey).(*config.Config),
	}
	return c.JSON(http.StatusOK, response)
}
