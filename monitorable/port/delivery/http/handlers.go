package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/port"
	portModels "github.com/monitoror/monitoror/monitorable/port/models"

	"github.com/labstack/echo/v4"
)

type PortDelivery struct {
	portUsecase port.Usecase
}

func NewPortDelivery(p port.Usecase) *PortDelivery {
	return &PortDelivery{p}
}

func (h *PortDelivery) GetPort(c echo.Context) error {
	// Bind / check Params
	params := &portModels.PortParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.portUsecase.Port(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
