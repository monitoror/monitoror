package info

import (
	"net/http"

	"github.com/jsdidierlaurent/monitowall/config"

	"github.com/jsdidierlaurent/monitowall/errors"

	"github.com/jsdidierlaurent/monitowall/middlewares"
	"github.com/jsdidierlaurent/monitowall/renderings"

	"github.com/labstack/echo/v4"
)

func GetInfo(c echo.Context) error {
	value := c.Get(middlewares.BuildInfoContextKey)
	if value == nil {
		return errors.NewSystemError("unable to find config.BuildInfo in echo.context")
	}
	buildInfo, ok := value.(*config.BuildInfo)
	if !ok {
		return errors.NewSystemError("unable cast value in *config.BuildInfo")
	}

	value = c.Get(middlewares.ConfigContextKey)
	if value == nil {
		return errors.NewSystemError("unable to find config.Config in echo.context")
	}
	conf, ok := value.(*config.Config)
	if !ok {
		return errors.NewSystemError("unable cast value in *config.Config")
	}

	response := renderings.InfoResponse{
		BuildInfo: *buildInfo,
		Config:    *conf,
	}
	return c.JSON(http.StatusOK, response)
}
