package models

import "github.com/labstack/echo/v4"

type (
	MonitororError interface {
		error
		Send(ctx echo.Context)
	}
)
