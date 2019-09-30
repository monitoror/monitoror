package http

import (
	"net/http"

	"github.com/monitoror/monitoror/monitorable/pingdom"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/pingdom/models"

	"github.com/labstack/echo/v4"
)

type httpPingdomDelivery struct {
	pingdomUsecase pingdom.Usecase
}

func NewHttpPingdomDelivery(p pingdom.Usecase) *httpPingdomDelivery {
	return &httpPingdomDelivery{p}
}

func (h *httpPingdomDelivery) GetCheck(c echo.Context) error {
	// Bind / Check Params
	params := &models.CheckParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return QueryParamsError
	}

	tile, err := h.pingdomUsecase.Check(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
