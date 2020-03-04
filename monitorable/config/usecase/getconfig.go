package usecase

import (
	"fmt"

	"github.com/monitoror/monitoror/monitorable/config/models"
)

// GetConfig and set default value for Config from repository
func (cu *configUsecase) GetConfig(params *models.ConfigParams) *models.ConfigBag {
	configBag := &models.ConfigBag{}

	var err error
	if params.URL != "" {
		configBag.Config, err = cu.repository.GetConfigFromURL(params.URL)
	} else if params.Path != "" {
		configBag.Config, err = cu.repository.GetConfigFromPath(params.Path)
	}

	if err == nil {
		return configBag
	}

	switch e := err.(type) {
	case *models.ConfigFileNotFoundError:
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorConfigNotFound,
			Message: e.Error(),
			Data:    models.ConfigErrorData{Value: e.PathOrURL},
		})
	case *models.ConfigVersionFormatError:
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnsupportedVersion,
			Message: e.Error(),
			Data: models.ConfigErrorData{
				FieldName: "version",
				Value:     e.WrongVersion,
				Expected:  fmt.Sprintf(`"%s" >= version >= "%s"`, MinimalVersion, CurrentVersion),
			},
		})
	case *models.ConfigUnmarshalError:
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnableToParseConfig,
			Message: e.Error(),
			Data: models.ConfigErrorData{
				ConfigExtract: e.RawConfig,
			},
		})
	default:
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnexpectedError,
			Message: err.Error(),
		})
	}

	return configBag
}
