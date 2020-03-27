package http

import (
	"net/http"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/api"
	"github.com/monitoror/monitoror/monitorables/jenkins/api/models"

	"github.com/labstack/echo/v4"
)

type JenkinsDelivery struct {
	jenkinsUsecase api.Usecase
}

func NewJenkinsDelivery(p api.Usecase) *JenkinsDelivery {
	return &JenkinsDelivery{p}
}

func (h *JenkinsDelivery) GetBuild(c echo.Context) error {
	// Bind / check Params
	params := &models.BuildParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return coreModels.QueryParamsError
	}

	tile, err := h.jenkinsUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
