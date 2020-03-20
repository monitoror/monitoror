package http

import (
	"net/http"

	coreModels "github.com/monitoror/monitoror/models"
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
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return coreModels.QueryParamsError
	}

	tile, err := h.travisciUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
