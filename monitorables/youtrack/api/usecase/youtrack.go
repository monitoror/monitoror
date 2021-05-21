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
	tile.Label = "Youtrack issues count"

	issues, err := yu.repository.GetIssues(params.Query)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to find build"}
	}
	issuesCount := len(*issues)

	tile.Status = coreModels.SuccessStatus
	tile.Metrics.Values = append(tile.Metrics.Values, fmt.Sprintf("%d", issuesCount))

	// Count threshold
	if len(params.CountThreshold) != 0 {
		tile.Status = params.CountThreshold.GetTileStatus(issuesCount, tile.Status)
	}

	// Priority threshold
	if len(params.PriorityFieldThreshold) != 0 {
		priorityIssueCount := 0
		for _, issue := range *issues {
			for _, field := range issue.CustomFields {
				if field.Name == params.GetPriorityFieldLabelWithDefault() && field.Value != nil && field.Value.(map[string]interface{})["name"] == params.GetPriorityFieldValueWithDefault() {
					priorityIssueCount++
				}
			}
		}

		tile.Status = params.PriorityFieldThreshold.GetTileStatus(priorityIssueCount, tile.Status)
	}

	return tile, nil
}
