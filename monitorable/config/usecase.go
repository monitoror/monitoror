package config

import (
	"github.com/monitoror/monitoror/monitorable/config/models"
)

// Usecase represent the config's usecases
type (
	Usecase interface {
		Config(params *models.ConfigParams) (*models.Config, error)
		Verify(config *models.Config) error
		Hydrate(config *models.Config) error
	}
)
