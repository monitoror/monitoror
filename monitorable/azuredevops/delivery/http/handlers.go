package http

import (
	"net/http"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	"github.com/monitoror/monitoror/monitorable/azuredevops/models"

	"github.com/labstack/echo/v4"
)

type httpJenkinsDelivery struct {
	azureDevOpsUsecase azuredevops.Usecase
}

func NewHttpAzureDevOpsDelivery(p azuredevops.Usecase) *httpJenkinsDelivery {
	return &httpJenkinsDelivery{p}
}

func (h *httpJenkinsDelivery) GetBuild(c echo.Context) error {
	// Bind / check Params
	params := &models.BuildParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return QueryParamsError
	}

	tile, err := h.azureDevOpsUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *httpJenkinsDelivery) GetRelease(c echo.Context) error {
	// Bind / check Params
	params := &models.ReleaseParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return QueryParamsError
	}

	tile, err := h.azureDevOpsUsecase.Release(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
