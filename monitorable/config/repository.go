package config

import (
	"github.com/monitoror/monitoror/monitorable/config/models"
)

type (
	Repository interface {
		GetConfigFromURL(string) (*models.Config, error)
		GetConfigFromPath(string) (*models.Config, error)
	}
)
