//+build !faker

package usecase

import (
	"github.com/monitoror/monitoror/monitorable/config/models"
)

// GetConfig and set default value for Config from repository
func (cu *configUsecase) GetConfig(params *models.ConfigParams) (config *models.Config, err error) {
	if params.Url != "" {
		config, err = cu.repository.GetConfigFromUrl(params.Url)
	} else if params.Path != "" {
		config, err = cu.repository.GetConfigFromPath(params.Path)
	}

	if err != nil {
		return
	}

	// Set config to CurrentVersion if config isn't set
	if config.Version == 0 {
		config.Version = CurrentVersion
	}

	// Clean Errors / Warnings
	config.Errors = []string{}
	config.Warnings = []string{}

	return
}
