package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/gitea/api"
	"github.com/monitoror/monitoror/monitorables/gitea/api/models"

	"github.com/labstack/echo/v4"
)

type GiteaDelivery struct {
	giteaUsecase api.Usecase
}

func NewGiteaDelivery(p api.Usecase) *GiteaDelivery {
	return &GiteaDelivery{p}
}

func (h *GiteaDelivery) GetCount(c echo.Context) error {
	// Bind / check Params
	params := &models.CountParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.giteaUsecase.Count(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *GiteaDelivery) GetChecks(c echo.Context) error {
	// Bind / check Params
	params := &models.ChecksParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.giteaUsecase.Checks(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *GiteaDelivery) GetPullRequest(c echo.Context) error {
	// Bind / check Params
	params := &models.PullRequestParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := h.giteaUsecase.PullRequest(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
