package config

import (
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

// Usecase represent the config's usecases
type (
	Helper interface {
		RegisterTile(tileType tiles.TileType, path string, validator utils.Validator)
	}

	Usecase interface {
		Helper

		GetConfig(params *models.ConfigParams) (*models.Config, error)
		Verify(config *models.Config) error
		Hydrate(config *models.Config, host string) error
	}
)
