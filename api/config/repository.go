package config

import (
	"github.com/monitoror/monitoror/api/config/models"
)

type (
	Repository interface {
		GetConfigFromURL(string) (*models.Config, error)
		GetConfigFromPath(string) (*models.Config, error)
	}
)
