package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
)

const (
	GithubCountTileType       coreModels.TileType = "GITHUB-COUNT"
	GithubChecksTileType      coreModels.TileType = "GITHUB-CHECKS"
	GithubPullRequestTileType coreModels.TileType = "GITHUB-PULLREQUESTS"
)

type (
	Usecase interface {
		Count(params *models.CountParams) (*coreModels.Tile, error)
		Checks(params *models.ChecksParams) (*coreModels.Tile, error)
		PullRequests(params interface{}) ([]builder.Result, error)
	}
)
