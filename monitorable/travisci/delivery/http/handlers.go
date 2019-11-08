package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/travisci"
	travisCIModels "github.com/monitoror/monitoror/monitorable/travisci/models"

	"github.com/labstack/echo/v4"
)

type TravisCIDelivery struct {
	travisciUsecase travisci.Usecase
}

func NewTravisCIDelivery(p travisci.Usecase) *TravisCIDelivery {
	return &TravisCIDelivery{p}
}

func (h *TravisCIDelivery) GetBuild(c echo.Context) error {
	// Bind / check Params
	params := &travisCIModels.BuildParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.travisciUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
