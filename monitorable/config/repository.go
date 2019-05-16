package config

import "github.com/monitoror/monitoror/monitorable/config/models"

// Repository represent the config's repository contract
type (
	Repository interface {
		GetConfigFromUrl(string) (*models.Config, error)
		GetConfigFromPath(string) (*models.Config, error)
	}
)
