package config

import (
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config/models"
	. "github.com/monitoror/monitoror/pkg/monitoror/validator"
)

// Usecase represent the config's usecases
type (
	Helper interface {
		RegisterTile(tileType tiles.TileType, path string, validator Validator)
		RegisterTileWithConfigVariant(tileType tiles.TileType, configVariant, path string, validator Validator)
	}

	Usecase interface {
		Helper

		GetConfig(params *models.ConfigParams) (*models.Config, error)
		Verify(config *models.Config) error
		Hydrate(config *models.Config, host string) error
	}
)
