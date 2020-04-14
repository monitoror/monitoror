package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
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
	if err := delivery.BindAndValidateRequestParams(c, params); err != nil {
		return err
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
	if err := delivery.BindAndValidateRequestParams(c, params); err != nil {
		return err
	}

	tile, err := h.githubUsecase.Checks(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
