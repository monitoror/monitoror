package models

import (
	"regexp"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/slice"
)

type (
	StatusCodesProvider interface {
		GetStatusCodes() (min int, max int)
	}

	RegexProvider interface {
		GetRegex() string
		GetRegexp() *regexp.Regexp
	}

	FormatedDataProvider interface {
		GetFormat() string
		GetKey() string
	}
)

const (
	DefaultMinStatusCode = 200
	DefaultMaxStatusCode = 399
)

const (
	JSONFormat = "JSON"
	YAMLFormat = "YAML"
	XMLFormat  = "XML"
)

var availableFormats = []string{JSONFormat, YAMLFormat, XMLFormat}

func isValid(url string, statusCodesProvider StatusCodesProvider) bool {
	if url == "" {
		return false
	}

	min, max := statusCodesProvider.GetStatusCodes()
	return min <= max
}

func isValidRegex(regexProvider RegexProvider) bool {
	regex := regexProvider.GetRegex()
	if regex != "" {
		_, err := regexp.Compile(regex)
		if err != nil {
			return false
		}
	}

	return true
}

func isValidKey(formatedDataProvider FormatedDataProvider) bool {
	key := formatedDataProvider.GetKey()
	if key == "" || key == "." {
		return false
	}

	return true
}

func isValidFormat(formatedDataProvider FormatedDataProvider) bool {
	format := formatedDataProvider.GetFormat()
	if _, find := slice.Find(availableFormats, format); !find {
		return false
	}

	return true
}

func getStatusCodes(statusCodeMin, statusCodeMax *int) (min int, max int) {
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
		return regexp.MustCompile(regex) // Already validate by isValid
	}
	return nil
}
