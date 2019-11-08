package travisci

import (
	"github.com/monitoror/monitoror/models"
	travisCIModels "github.com/monitoror/monitoror/monitorable/travisci/models"
)

const (
	TravisCIBuildTileType models.TileType = "TRAVISCI-BUILD"
)

type (
	Usecase interface {
		Build(params *travisCIModels.BuildParams) (*models.Tile, error)
	}
)
