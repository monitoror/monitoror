package errors

import (
	"fmt"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	"github.com/monitoror/monitoror/models/tiles"
)

type QueryParamsError struct {
	err error
}

func NewQueryParamsError(err error) *QueryParamsError {
	return &QueryParamsError{err}
}

func (qpe *QueryParamsError) Send(ctx echo.Context) {
	log.Warn(qpe.Error())
	_ = ctx.JSON(http.StatusBadRequest, tiles.NewErrorTile("Invalid configuration", qpe.Error()))
}

func (qpe *QueryParamsError) Error() (err string) {
	err = fmt.Sprintf("Unable to parse/check queryParams into struct")
	if qpe.err != nil {
		err += fmt.Sprintf(", %v", qpe.err)
	}
	return
}
