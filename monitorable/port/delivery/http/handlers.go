package http

import (
	"net/http"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/port"
	"github.com/monitoror/monitoror/monitorable/port/models"

	"github.com/labstack/echo/v4"
)

type httpPortDelivery struct {
	portUsecase port.Usecase
}

func NewHttpPortDelivery(p port.Usecase) *httpPortDelivery {
	return &httpPortDelivery{p}
}

func (h *httpPortDelivery) GetPort(c echo.Context) error {
	// Bind / check Params
	params := &models.PortParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return QueryParamsError
	}

	tile, err := h.portUsecase.Port(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
