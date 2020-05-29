package http

import (
	"net/http"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/delivery"
	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"

	"github.com/labstack/echo/v4"
)

type GitlabDelivery struct {
	gitlabUsecase api.Usecase
}

func NewGitlabDelivery(p api.Usecase) *GitlabDelivery {
	return &GitlabDelivery{p}
}

func (gd *GitlabDelivery) GetIssues(c echo.Context) error {
	// Bind / check Params
	params := &models.IssuesParams{}
	_ = delivery.BindAndValidateParams(c, params)

	tile, err := gd.gitlabUsecase.Issues(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (gd *GitlabDelivery) GetPipeline(c echo.Context) error {
	// Bind / check Params
	params := &models.PipelineParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := gd.gitlabUsecase.Pipeline(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (gd *GitlabDelivery) GetMergeRequest(c echo.Context) error {
	// Bind / check Params
	params := &models.MergeRequestParams{}
	if err := delivery.BindAndValidateParams(c, params); err != nil {
		return err
	}

	tile, err := gd.gitlabUsecase.MergeRequest(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
