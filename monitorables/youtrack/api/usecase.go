package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/models"
)

const (
	YoutrackCountIssuesTileType coreModels.TileType = "YOUTRACK-COUNT-ISSUES"
)

type (
	Usecase interface {
		CountIssues(params *models.IssuesCountParams) (*coreModels.Tile, error)
	}
)
