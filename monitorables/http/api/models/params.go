package models

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	pkgConfig "github.com/monitoror/monitoror/internal/pkg/api/config"
)

type (
	GenericParamsProvider interface {
		GetURL() (url string)
		GetStatusCodes() (min int, max int)
	}

	RegexParamsProvider interface {
		GetRegex() string
		GetRegexp() *regexp.Regexp
	}

	FormattedParamsProvider interface {
		GetFormat() Format
		GetKey() string
	}

	Format string
)

const (
	DefaultMinStatusCode = 200
	DefaultMaxStatusCode = 399
)

const (
	JSONFormat Format = "JSON"
	YAMLFormat Format = "YAML"
	XMLFormat  Format = "XML"
)

var supportedFormats = map[Format]bool{
	JSONFormat: true,
	YAMLFormat: true,
	XMLFormat:  true,
}

func validateURL(params GenericParamsProvider) *uiConfigModels.ConfigError {
	u := params.GetURL()
	if u == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "url" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "url"},
		}
	}

	if _, err := url.Parse(u); err != nil {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorInvalidFieldValue,
			Message: fmt.Sprintf(`Invalid "url" field. Must be a valid URL.`),
			Data: uiConfigModels.ConfigErrorData{
				FieldName: "url",
				Expected:  "valid URL",
			},
		}
	}

	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorInvalidFieldValue,
			Message: fmt.Sprintf(`Unsupported "url" protocol. Must be http or https`),
			Data: uiConfigModels.ConfigErrorData{
				FieldName: "url",
				Expected:  "http or https protocol",
			},
		}
	}

	return nil
}

func validateStatusCode(params GenericParamsProvider) *uiConfigModels.ConfigError {
	if min, max := params.GetStatusCodes(); min > max {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorInvalidFieldValue,
			Message: fmt.Sprintf(`Invalid "statusCodeMin" field. Must be lower or equal to statusCodeMax.`),
			Data: uiConfigModels.ConfigErrorData{
				FieldName: "statusCodeMin",
				Expected:  "statusCodeMin <= statusCodeMax",
			},
		}
	}

	return nil
}

func validateRegex(params RegexParamsProvider) *uiConfigModels.ConfigError {
	regex := params.GetRegex()
	if regex != "" {
		_, err := regexp.Compile(regex)
		if err != nil {
			return &uiConfigModels.ConfigError{
				ID:      uiConfigModels.ConfigErrorInvalidFieldValue,
				Message: fmt.Sprintf(`Invalid "regex" field. Must be a valid golang regex.`),
				Data: uiConfigModels.ConfigErrorData{
					FieldName: "regex",
					Expected:  "valid golang regex",
				},
			}
		}
	}

	return nil
}

func validateKey(params FormattedParamsProvider) *uiConfigModels.ConfigError {
	key := params.GetKey()
	if key == "" || key == "." {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "key" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "key"},
		}
	}

	return nil
}

func validateFormat(params FormattedParamsProvider) *uiConfigModels.ConfigError {
	format := params.GetFormat()
	if format == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "format" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "format"},
		}
	}

	if find := supportedFormats[format]; !find {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorInvalidFieldValue,
			Message: fmt.Sprintf(`Unknown %q format in tile http formatted params. Must be %s`, format, pkgConfig.Keys(supportedFormats)),
			Data: uiConfigModels.ConfigErrorData{
				FieldName: "format",
				Value:     string(format),
				Expected:  pkgConfig.Keys(supportedFormats),
			},
		}
	}

	return nil
}

func getStatusCodesWithDefault(statusCodeMin, statusCodeMax *int) (min int, max int) {
	min = DefaultMinStatusCode
	if statusCodeMin != nil {
		min = *statusCodeMin
	}
	max = DefaultMaxStatusCode
	if statusCodeMax != nil {
		max = *statusCodeMax
	}
	return
}

func getRegexp(regex string) *regexp.Regexp {
	if regex != "" {
		return regexp.MustCompile(regex) // Already validate by validateRegex
	}
	return nil
}
