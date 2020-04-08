package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
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
	if err := delivery.BindAndValidateRequestParams(c, params); err != nil {
		return err
	}

	tile, err := h.jenkinsUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
