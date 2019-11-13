package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	azureModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"

	"github.com/labstack/echo/v4"
)

type AzureDevOpsDelivery struct {
	azureDevOpsUsecase azuredevops.Usecase
}

func NewAzureDevOpsDelivery(p azuredevops.Usecase) *AzureDevOpsDelivery {
	return &AzureDevOpsDelivery{p}
}

func (h *AzureDevOpsDelivery) GetBuild(c echo.Context) error {
	// Bind / check Params
	params := &azureModels.BuildParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.azureDevOpsUsecase.Build(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *AzureDevOpsDelivery) GetRelease(c echo.Context) error {
	// Bind / check Params
	params := &azureModels.ReleaseParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.azureDevOpsUsecase.Release(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
