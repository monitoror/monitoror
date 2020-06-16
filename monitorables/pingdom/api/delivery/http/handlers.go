package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
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
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.pingdomUsecase.Check(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *PingdomDelivery) GetTransactionCheck(c echo.Context) error {
	// Bind / Check Params
	params := &models.TransactionCheckParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.pingdomUsecase.TransactionCheck(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
