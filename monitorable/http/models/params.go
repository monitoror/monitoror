package models

import "regexp"

type (
	StatusCodesProvider interface {
		GetStatusCodes() (min int, max int)
	}

	RegexProvider interface {
		GetRegex() string
		GetRegexp() *regexp.Regexp
	}

	FormatedDataProvider interface {
		GetKey() string
		GetUnmarshaller() func(data []byte, v interface{}) error
	}
)

const (
	DefaultMinStatusCode = 200
	DefaultMaxStatusCode = 399
)

func isValid(url string, statusCodesProvider StatusCodesProvider) bool {
	if url == "" {
		return false
	}

	min, max := statusCodesProvider.GetStatusCodes()
	if min > max {
		return false
	}

	return true
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
