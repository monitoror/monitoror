package usecase

import (
	"errors"
	"fmt"

	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/monitorable/config/repository"
)

// GetConfig and set default value for Config from repository
func (cu *configUsecase) GetConfig(params *models.ConfigParams) (configBag *models.ConfigBag, err error) {
	configBag = &models.ConfigBag{}

	var config *models.Config
	if params.URL != "" {
		config, err = cu.repository.GetConfigFromURL(params.URL)
	} else if params.Path != "" {
		config, err = cu.repository.GetConfigFromPath(params.Path)
	}

	if err != nil {
		if errors.Is(err, repository.ErrConfigFileNotFound) {
			var pathOrURL string

			if params.URL != "" {
				pathOrURL = params.URL
			}
			if params.Path != "" {
				pathOrURL = params.Path
			}

			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorConfigNotFound,
				Message: fmt.Sprintf("Config not found at: %s", pathOrURL),
				Data: models.ConfigErrorData{
					Value: pathOrURL,
				},
			})
			err = nil
		}
	}

	if config != nil {
		configBag.Config = *config
	}

	return
}
