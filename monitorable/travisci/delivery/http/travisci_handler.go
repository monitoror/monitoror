package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/monitorable/travisci/model"

	"github.com/monitoror/monitoror/monitorable/travisci"

	"github.com/labstack/echo/v4"
)

type httpTravisCIHandler struct {
	travisciUsecase travisci.Usecase
}

func NewHttpTravisCIHandler(p travisci.Usecase) *httpTravisCIHandler {
	return &httpTravisCIHandler{p}
}

func (h *httpTravisCIHandler) GetTravisCIBuild(c echo.Context) error {
	// Bind / check Params
	params := &model.BuildParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return errors.NewQueryParamsError(err)
	}

	tile, err := h.travisciUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
