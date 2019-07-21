package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models/errors"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	"github.com/monitoror/monitoror/monitorable/jenkins/models"

	"github.com/labstack/echo/v4"
)

type httpJenkinsDelivery struct {
	jenkinsUsecase jenkins.Usecase
}

func NewHttpJenkinsDelivery(p jenkins.Usecase) *httpJenkinsDelivery {
	return &httpJenkinsDelivery{p}
}

func (h *httpJenkinsDelivery) GetBuild(c echo.Context) error {
	// Bind / check Params
	params := &models.BuildParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return errors.NewQueryParamsError(err)
	}

	tile, err := h.jenkinsUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
