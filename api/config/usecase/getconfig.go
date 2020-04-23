package usecase

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/api/config/versions"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/internal/pkg/api/config"
	"github.com/monitoror/monitoror/internal/pkg/validator/validate"

	"github.com/fatih/structs"
)

var (
	unknownFieldRegex            *regexp.Regexp
	fieldTypeMismatchRegex       *regexp.Regexp
	invalidEscapedCharacterRegex *regexp.Regexp

	urlRegex *regexp.Regexp
)

func init() {
	// Based on: https://github.com/golang/go/blob/release-branch.go1.14/src/encoding/json/decode.go#L755
	unknownFieldRegex = regexp.MustCompile(`json: unknown field "(.*)"`)

	// Based on: https://github.com/golang/go/blob/go1.14/src/encoding/json/decode.go#L134
	fieldTypeMismatchRegex = regexp.MustCompile(`json: cannot unmarshal .+ into Go struct field (.+) of type (.+)`)

	// Based on: https://github.com/golang/go/blob/go1.14/src/encoding/json/scanner.go#L343
	invalidEscapedCharacterRegex = regexp.MustCompile(`'(.*)' in string escape code`)

	// Simple regex that identifies url
	urlRegex = regexp.MustCompile(validate.HTTPRegex)
}

// GetConfig and set default value for Config from repository
func (cu *configUsecase) GetConfig(params *models.ConfigParams) *models.ConfigBag {
	configBag := &models.ConfigBag{}
	var err error

	// Lookup for a url
	if urlRegex.MatchString(params.Config) {
		configBag.Config, err = cu.repository.GetConfigFromURL(params.Config)
	} else {
		// If config is empty, set default value
		configName := coreConfig.ConfigName(params.Config)
		if configName == "" {
			configName = coreConfig.DefaultConfigName
		}

		// Lookup for a named Config
		if namedConfig, ok := cu.namedConfigs[configName]; ok {
			if urlRegex.MatchString(namedConfig) {
				configBag.Config, err = cu.repository.GetConfigFromURL(namedConfig)
			} else {
				configBag.Config, err = cu.repository.GetConfigFromPath(namedConfig)
			}
		} else {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnknownNamedConfig,
				Message: fmt.Sprintf(`Unknown %q named config. Must be %s`, params.Config, config.Keys(cu.namedConfigs)),
				Data: models.ConfigErrorData{
					Value:    params.Config,
					Expected: config.Keys(cu.namedConfigs),
				},
			})
		}
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
	case *versions.ConfigVersionFormatError:
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnsupportedVersion,
			Message: e.Error(),
			Data: models.ConfigErrorData{
				FieldName: "version",
				Value:     e.WrongVersion,
				Expected:  fmt.Sprintf(`%q >= version >= %q`, versions.MinimalVersion, versions.CurrentVersion),
			},
		})
	case *models.ConfigUnmarshalError:
		// Check if error is "json: unknown field"
		if unknownFieldRegex.MatchString(err.Error()) {
			subMatch := unknownFieldRegex.FindAllStringSubmatch(err.Error(), 1)

			var field string
			if len(subMatch) > 0 && len(subMatch[0]) > 1 {
				field = subMatch[0][1]
			}

			configField := structs.Fields(models.Config{})
			tileConfigFields := structs.Fields(models.TileConfig{})
			expectedFields := append(configField, tileConfigFields...)
			var expectedFieldNames []string

			for _, expectedField := range expectedFields {
				jsonTag := expectedField.Tag("json")
				if jsonTag != "" && jsonTag != "-" {
					expectedFieldName := strings.Replace(jsonTag, ",omitempty", "", 1)
					expectedFieldNames = append(expectedFieldNames, expectedFieldName)
				}
			}

			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnknownField,
				Message: e.Error(),
				Data: models.ConfigErrorData{
					FieldName:     field,
					ConfigExtract: e.RawConfig,
					Expected:      strings.Join(expectedFieldNames, ", "),
				},
			})
		} else if fieldTypeMismatchRegex.MatchString(err.Error()) {
			subMatch := fieldTypeMismatchRegex.FindAllStringSubmatch(err.Error(), 1)

			var field string
			var expectedType string
			if len(subMatch) > 0 && len(subMatch[0]) > 1 {
				fieldParts := strings.Split(subMatch[0][1], ".")
				field = fieldParts[len(fieldParts)-1]
				expectedType = subMatch[0][2]
			}

			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorFieldTypeMismatch,
				Message: e.Error(),
				Data: models.ConfigErrorData{
					FieldName:     field,
					ConfigExtract: e.RawConfig,
					Expected:      expectedType,
				},
			})
		} else if invalidEscapedCharacterRegex.MatchString(err.Error()) {
			subMatch := invalidEscapedCharacterRegex.FindAllStringSubmatch(err.Error(), 1)

			var invalidEscapedCharacter string
			if len(subMatch) > 0 && len(subMatch[0]) > 1 {
				invalidEscapedCharacter = subMatch[0][1]
			}

			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorInvalidEscapedCharacter,
				Message: e.Error(),
				Data: models.ConfigErrorData{
					ConfigExtract:          e.RawConfig,
					ConfigExtractHighlight: fmt.Sprintf(`\%s`, invalidEscapedCharacter),
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
