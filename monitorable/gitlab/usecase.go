package gitlab

import (
	"github.com/monitoror/monitoror/models"
	gitlabModels "github.com/monitoror/monitoror/monitorable/gitlab/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
)

const (
	GitlabCountTileType       models.TileType = "GITLAB-COUNT"
	GitlabPipelinesTileType   models.TileType = "GITLAB-PIPELINES"
	GitlabPullRequestTileType models.TileType = "GITLAB-MERGEREQUESTS"
)

type (
	Usecase interface {
		Count(params *gitlabModels.CountParams) (*models.Tile, error)
		Pipelines(params *gitlabModels.PipelinesParams) (*models.Tile, error)
		ListDynamicTile(params interface{}) ([]builder.Result, error)
	}
)
