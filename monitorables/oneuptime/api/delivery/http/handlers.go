package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api/models"

	"github.com/labstack/echo/v4"
)

type OneuptimeDelivery struct {
	oneuptimeUsecase api.Usecase
}

func NewOneuptimeDelivery(p api.Usecase) *OneuptimeDelivery {
	return &OneuptimeDelivery{p}
}

func (h *OneuptimeDelivery) GetMonitor(c echo.Context) error {
	// Bind / Check Params
	params := &models.MonitorParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.oneuptimeUsecase.Monitor(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
