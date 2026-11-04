//go:build faker

package usecase

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/youtrack/api"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/models"
	"github.com/monitoror/monitoror/pkg/nonempty"
)

type (
	youtrackUsecase struct{}
)

func NewYoutrackUsecase() api.Usecase {
	return &youtrackUsecase{}
}

func (yu *youtrackUsecase) CountIssues(params *models.IssuesCountParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.YoutrackCountIssuesTileType).WithMetrics(coreModels.NumberUnit)
	tile.Label = "YouTrack issues"

	tile.Status = nonempty.Struct(params.Status, coreModels.SuccessStatus).(coreModels.TileStatus)

	if len(params.ValueValues) != 0 {
		tile.Metrics.Values = params.ValueValues
	} else {
		tile.Metrics.Values = []string{"42"}
	}

	return tile, nil
}
