//+build !faker

package usecase

import (
	"github.com/monitoror/monitoror/monitorable/config/models"
)

// GetConfig and set default value for Config from repository
func (cu *configUsecase) GetConfig(params *models.ConfigParams) (config *models.Config, err error) {
	if params.URL != "" {
		config, err = cu.repository.GetConfigFromURL(params.URL)
	} else if params.Path != "" {
		config, err = cu.repository.GetConfigFromPath(params.Path)
	}

	if err != nil {
		return
	}

	// Clean Errors / Warnings
	config.Errors = []string{}
	config.Warnings = []string{}

	return
}
