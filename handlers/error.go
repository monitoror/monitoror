package handlers

import (
	"net/http"

	"github.com/jsdidierlaurent/monitowall/errors"
	"github.com/labstack/echo/v4"
)

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func JSONErrorHandler(err error, c echo.Context) {

	switch e := err.(type) {
	case *errors.SystemError:
		// System Error
		_ = c.JSON(http.StatusOK, ApiError{
			Status:  500,
			Message: e.Message,
		})
	default:
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == 404 {
				// 404
				_ = c.JSON(he.Code, ApiError{
					Status:  he.Code,
					Message: "Not Found",
				})
			} else {
				_ = c.JSON(he.Code, ApiError{
					Status:  he.Code,
					Message: "System Error",
				})
			}
		}
	}

}
