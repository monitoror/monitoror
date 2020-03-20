package http

import (
	"net/http"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api"
	"github.com/monitoror/monitoror/monitorables/github/api/models"

	"github.com/labstack/echo/v4"
)

type GithubDelivery struct {
	githubUsecase api.Usecase
}

func NewGithubDelivery(p api.Usecase) *GithubDelivery {
	return &GithubDelivery{p}
}

func (h *GithubDelivery) GetCount(c echo.Context) error {
	// Bind / check Params
	params := &models.CountParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return coreModels.QueryParamsError
	}

	tile, err := h.githubUsecase.Count(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *GithubDelivery) GetChecks(c echo.Context) error {
	// Bind / check Params
	params := &models.ChecksParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return coreModels.QueryParamsError
	}

	tile, err := h.githubUsecase.Checks(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
