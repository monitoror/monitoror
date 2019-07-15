package travisci

import (
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/travisci/models"
)

const (
	TravisCIBuildTileType TileType = "TRAVISCI-BUILD"
)

// Usecase represent the travisci's usecases
type (
	Usecase interface {
		Build(params *models.BuildParams) (*BuildTile, error)
	}
)
