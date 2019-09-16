package http

import (
	"net/http"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/ping"
	"github.com/monitoror/monitoror/monitorable/ping/models"

	"github.com/labstack/echo/v4"
)

type httpPingDelivery struct {
	pingUsecase ping.Usecase
}

func NewHttpPingDelivery(p ping.Usecase) *httpPingDelivery {
	return &httpPingDelivery{p}
}

func (h *httpPingDelivery) GetPing(c echo.Context) error {
	// Bind / Check Params
	params := &models.PingParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return QueryParamsError
	}

	tile, err := h.pingUsecase.Ping(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
