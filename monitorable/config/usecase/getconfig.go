package usecase

import (
	"fmt"
	"regexp"

	"github.com/monitoror/monitoror/monitorable/config/models"
)

var unknownFieldRegex *regexp.Regexp

func init() {
	// Based on: https://github.com/golang/go/blob/release-branch.go1.14/src/encoding/json/decode.go#L755
	unknownFieldRegex = regexp.MustCompile(`json: unknown field "(.*)"`)
}

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
		// Check if error is "json: unknown field"
		if unknownFieldRegex.MatchString(err.Error()) {
			subMatch := unknownFieldRegex.FindAllStringSubmatch(err.Error(), 1)

			var field = ""
			if len(subMatch) > 0 && len(subMatch[0]) > 1 {
				field = subMatch[0][1]
			}

			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnknownField,
				Message: e.Error(),
				Data: models.ConfigErrorData{
					FieldName:     field,
					ConfigExtract: e.RawConfig,
				},
			})
		} else {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnableToParseConfig,
				Message: e.Error(),
				Data: models.ConfigErrorData{
					ConfigExtract: e.RawConfig,
				},
			})
		}
	default:
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnexpectedError,
			Message: err.Error(),
		})
	}

	return configBag
}
