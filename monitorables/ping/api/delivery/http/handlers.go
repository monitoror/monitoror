package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	"github.com/monitoror/monitoror/monitorables/ping/api/models"

	"github.com/labstack/echo/v4"
)

type PingDelivery struct {
	pingUsecase api.Usecase
}

func NewPingDelivery(p api.Usecase) *PingDelivery {
	return &PingDelivery{p}
}

func (h *PingDelivery) GetPing(c echo.Context) error {
	// Bind / Check Params
	params := &models.PingParams{}
	if err := delivery.BindAndValidateRequestParams(c, params); err != nil {
		return err
	}

	tile, err := h.pingUsecase.Ping(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
