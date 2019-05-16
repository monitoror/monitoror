package errors

import (
	"fmt"
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	"github.com/monitoror/monitoror/models/tiles"
)

type SystemError struct {
	message string
	err     error
}

func NewSystemError(message string, err error) *SystemError {
	return &SystemError{message, err}
}

func (se *SystemError) Send(ctx echo.Context) {
	log.Error(se.Error())
	_ = ctx.JSON(http.StatusInternalServerError, tiles.NewErrorTile("System Error", se.Error()))
}

func (se *SystemError) Error() (err string) {
	err = se.message
	if se.err != nil {
		err += fmt.Sprintf(", %v", se.err)
	}
	return
}
