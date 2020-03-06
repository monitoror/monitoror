package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/gitlab"
	gitlabModels "github.com/monitoror/monitoror/monitorable/gitlab/models"

	"github.com/labstack/echo/v4"
)

type GitlabDelivery struct {
	gitlabUsecase gitlab.Usecase
}

func NewGitlabDelivery(p gitlab.Usecase) *GitlabDelivery {
	return &GitlabDelivery{p}
}

func (h *GitlabDelivery) GetCount(c echo.Context) error {
	// Bind / check Params
	params := &gitlabModels.CountParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.gitlabUsecase.Count(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *GitlabDelivery) GetPipelines(c echo.Context) error {
	// Bind / check Params
	params := &gitlabModels.PipelinesParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.gitlabUsecase.Pipelines(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
