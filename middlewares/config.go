package middlewares

import (
	"github.com/jsdidierlaurent/monitowall/config"
	"github.com/labstack/echo/v4"
)

const (
	ConfigContextKey    = "jsdidierlaurent.monitowall.config"
	BuildInfoContextKey = "jsdidierlaurent.monitowall.buildInfo"
)

//ConfigMiddleware Provide config and buildInfo to all route using echo.context
func ConfigMiddleware(config *config.Config, buildInfo *config.BuildInfo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Inject config in echo context
			c.Set(ConfigContextKey, config)
			// Inject buildInfo in echo context
			c.Set(BuildInfoContextKey, buildInfo)
			return next(c)
		}
	}
}
