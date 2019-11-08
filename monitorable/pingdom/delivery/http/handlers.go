package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	pingdomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"

	"github.com/labstack/echo/v4"
)

type PingdomDelivery struct {
	pingdomUsecase pingdom.Usecase
}

func NewPingdomDelivery(p pingdom.Usecase) *PingdomDelivery {
	return &PingdomDelivery{p}
}

func (h *PingdomDelivery) GetCheck(c echo.Context) error {
	// Bind / Check Params
	params := &pingdomModels.CheckParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.pingdomUsecase.Check(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
