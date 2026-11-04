//go:build !faker

package usecase

import (
	"fmt"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/youtrack/api"
	"github.com/monitoror/monitoror/monitorables/youtrack/api/models"
)

type (
	youtrackUsecase struct {
		repository api.Repository
	}
)

func NewYoutrackUsecase(repository api.Repository) api.Usecase {
	return &youtrackUsecase{repository}
}

func (yu *youtrackUsecase) CountIssues(params *models.IssuesCountParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.YoutrackCountIssuesTileType).WithMetrics(coreModels.NumberUnit)
	tile.Label = "YouTrack issues"

	count, err := yu.repository.CountIssues(params.Query)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to count issues"}
	}

	tile.Status = coreModels.SuccessStatus
	tile.Metrics.Values = append(tile.Metrics.Values, fmt.Sprintf("%d", count))

	if len(params.CountThreshold) != 0 {
		tile.Status = params.CountThreshold.GetTileStatus(count, tile.Status)
	}

	if len(params.PriorityFieldThreshold) != 0 {
		priorityCount, err := yu.repository.CountIssues(params.GetPriorityQuery())
		if err != nil {
			return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to count priority issues"}
		}

		priorityStatus := params.PriorityFieldThreshold.GetTileStatus(priorityCount, coreModels.SuccessStatus)
		tile.Status = coreModels.MostSevereTileStatus(tile.Status, priorityStatus)
	}

	return tile, nil
}
