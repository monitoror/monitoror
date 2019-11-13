package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/ping"
	pingModels "github.com/monitoror/monitoror/monitorable/ping/models"

	"github.com/labstack/echo/v4"
)

type PingDelivery struct {
	pingUsecase ping.Usecase
}

func NewPingDelivery(p ping.Usecase) *PingDelivery {
	return &PingDelivery{p}
}

func (h *PingDelivery) GetPing(c echo.Context) error {
	// Bind / Check Params
	params := &pingModels.PingParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.pingUsecase.Ping(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
