package http

import (
	"net/http"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api"
	"github.com/monitoror/monitoror/monitorables/port/api/models"

	"github.com/labstack/echo/v4"
)

type PortDelivery struct {
	portUsecase api.Usecase
}

func NewPortDelivery(p api.Usecase) *PortDelivery {
	return &PortDelivery{p}
}

func (h *PortDelivery) GetPort(c echo.Context) error {
	// Bind / check Params
	params := &models.PortParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return coreModels.QueryParamsError
	}

	tile, err := h.portUsecase.Port(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
