package github

import (
	"github.com/monitoror/monitoror/models"
	githubModels "github.com/monitoror/monitoror/monitorable/github/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
)

const (
	GithubIssuesTileType      models.TileType = "GITHUB-ISSUES"
	GithubChecksTileType      models.TileType = "GITHUB-CHECKS"
	GithubPullRequestTileType models.TileType = "GITHUB-PULLREQUESTS"
)

type (
	Usecase interface {
		Issues(params *githubModels.IssuesParams) (*models.Tile, error)
		Checks(params *githubModels.ChecksParams) (*models.Tile, error)
		ListDynamicTile(params interface{}) ([]builder.Result, error)
	}
)
