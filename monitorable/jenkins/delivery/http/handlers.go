package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"

	"github.com/labstack/echo/v4"
)

type JenkinsDelivery struct {
	jenkinsUsecase jenkins.Usecase
}

func NewJenkinsDelivery(p jenkins.Usecase) *JenkinsDelivery {
	return &JenkinsDelivery{p}
}

func (h *JenkinsDelivery) GetBuild(c echo.Context) error {
	// Bind / check Params
	params := &jenkinsModels.BuildParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.jenkinsUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
