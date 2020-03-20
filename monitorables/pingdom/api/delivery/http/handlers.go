package http

import (
	"net/http"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/pingdom/api"
	"github.com/monitoror/monitoror/monitorables/pingdom/api/models"

	"github.com/labstack/echo/v4"
)

type PingdomDelivery struct {
	pingdomUsecase api.Usecase
}

func NewPingdomDelivery(p api.Usecase) *PingdomDelivery {
	return &PingdomDelivery{p}
}

func (h *PingdomDelivery) GetCheck(c echo.Context) error {
	// Bind / Check Params
	params := &models.CheckParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return coreModels.QueryParamsError
	}

	tile, err := h.pingdomUsecase.Check(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
