package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/monitorable/port/models"

	"github.com/monitoror/monitoror/monitorable/port"

	"github.com/labstack/echo/v4"
)

type httpPortHandler struct {
	portUsecase port.Usecase
}

func NewHttpPortHandler(p port.Usecase) *httpPortHandler {
	return &httpPortHandler{p}
}

func (h *httpPortHandler) GetPort(c echo.Context) error {
	// Bind / check Params
	params := &models.PortParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return errors.NewQueryParamsError(err)
	}

	tile, err := h.portUsecase.Port(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
