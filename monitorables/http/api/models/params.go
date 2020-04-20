package models

import (
	"regexp"

	"github.com/monitoror/monitoror/internal/pkg/validator"
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

func validateStatusCode(params GenericParamsProvider) []validator.Error {
	if min, max := params.GetStatusCodes(); min > max {
		return []validator.Error{validator.NewDefaultError("StatusCodeMin", "statusCodeMin <= statusCodeMax")}
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
		r, _ := regexp.Compile(regex) // Already validate by validateRegex
		return r
	}
	return nil
}
