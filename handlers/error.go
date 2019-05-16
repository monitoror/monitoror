package handlers

import (
	"net/http"

	"github.com/monitoror/monitoror/models"

	"github.com/labstack/echo/v4"
)

type (
	ApiError struct {
		Code    int    `json:"status"`
		Message string `json:"message"`
	}
)

func HttpErrorHandler(err error, ctx echo.Context) {
	switch e := err.(type) {
	case models.MonitororError:
		e.Send(ctx)
	default:
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == 404 {
				// 404
				_ = ctx.JSON(he.Code, ApiError{
					Code:    he.Code,
					Message: "Not Found",
				})
				return
			}
		}

		_ = ctx.JSON(http.StatusInternalServerError, ApiError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
}
