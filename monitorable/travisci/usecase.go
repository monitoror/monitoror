package travisci

import (
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/travisci/model"
)

const (
	TravisCIBuildTileType TileType = "TRAVISCI-BUILD"
)

// Usecase represent the circleci's usecases
type (
	Usecase interface {
		Build(params *model.BuildParams) (*BuildTile, error)
	}
)
