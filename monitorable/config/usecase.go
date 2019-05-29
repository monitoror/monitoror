package config

import (
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

// Usecase represent the config's usecases
type (
	Regiterer interface {
		Register(tileType tiles.TileType, path string, configValidator utils.Validator)
	}

	Usecase interface {
		Regiterer

		Config(params *models.ConfigParams) (*models.Config, error)
		Verify(config *models.Config) error
		Hydrate(config *models.Config) error
	}
)
