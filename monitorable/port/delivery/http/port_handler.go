package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/monitorable/port/model"

	"github.com/monitoror/monitoror/monitorable/port"

	"github.com/labstack/echo/v4"
)

type HttpPortHandler struct {
	portUsecase port.Usecase
}

func NewHttpPortHandler(p port.Usecase) *HttpPortHandler {
	return &HttpPortHandler{p}
}

func (h *HttpPortHandler) GetPort(c echo.Context) error {
	// Bind / Validate Params
	params := &model.PortParams{}
	err := c.Bind(params)
	if err != nil || !params.Validate() {
		return errors.NewQueryParamsError(err)
	}

	tile, err := h.portUsecase.Port(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
