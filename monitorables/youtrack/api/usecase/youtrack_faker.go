//+build faker

package usecase

import (
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

	// TODO Seb

	return tile, nil
}
