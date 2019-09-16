package http

import (
	"net/http"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/travisci"
	"github.com/monitoror/monitoror/monitorable/travisci/models"

	"github.com/labstack/echo/v4"
)

type httpTravisCIDelivery struct {
	travisciUsecase travisci.Usecase
}

func NewHttpTravisCIDelivery(p travisci.Usecase) *httpTravisCIDelivery {
	return &httpTravisCIDelivery{p}
}

func (h *httpTravisCIDelivery) GetBuild(c echo.Context) error {
	// Bind / check Params
	params := &models.BuildParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return QueryParamsError
	}

	tile, err := h.travisciUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
