package http

import (
	"net/http"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api/models"

	"github.com/labstack/echo/v4"
)

type AzureDevOpsDelivery struct {
	azureDevOpsUsecase api.Usecase
}

func NewAzureDevOpsDelivery(p api.Usecase) *AzureDevOpsDelivery {
	return &AzureDevOpsDelivery{p}
}

func (h *AzureDevOpsDelivery) GetBuild(c echo.Context) error {
	// Bind / check Params
	params := &models.BuildParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return coreModels.QueryParamsError
	}

	tile, err := h.azureDevOpsUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *AzureDevOpsDelivery) GetRelease(c echo.Context) error {
	// Bind / check Params
	params := &models.ReleaseParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return coreModels.QueryParamsError
	}

	tile, err := h.azureDevOpsUsecase.Release(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
