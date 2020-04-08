package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/travisci/api"
	"github.com/monitoror/monitoror/monitorables/travisci/api/models"

	"github.com/labstack/echo/v4"
)

type TravisCIDelivery struct {
	travisciUsecase api.Usecase
}

func NewTravisCIDelivery(p api.Usecase) *TravisCIDelivery {
	return &TravisCIDelivery{p}
}

func (h *TravisCIDelivery) GetBuild(c echo.Context) error {
	// Bind / check Params
	params := &models.BuildParams{}
	if err := delivery.BindAndValidateRequestParams(c, params); err != nil {
		return err
	}

	tile, err := h.travisciUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
