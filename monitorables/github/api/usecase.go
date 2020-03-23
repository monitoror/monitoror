package api

import (
	models2 "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api/models"
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
		PullRequests(params interface{}) ([]models2.DynamicTileResult, error)
	}
)
