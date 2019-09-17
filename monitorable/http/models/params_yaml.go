//+build !faker

package models

import (
	"regexp"

	"gopkg.in/yaml.v2"
)

type (
	HttpYamlParams struct {
		Url           string `json:"url" query:"url"`
		Key           string `json:"key" query:"key"`
		Regex         string `json:"regex" query:"regex"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`
	}
)

func (p *HttpYamlParams) IsValid() bool {
	if !isValid(p.Url, p) {
		return false
	}

	if !isValidKey(p) {
		return false
	}

	return isValidRegex(p)
}

func (p *HttpYamlParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HttpYamlParams) GetRegex() string          { return p.Regex }
func (p *HttpYamlParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }

func (p *HttpYamlParams) GetKey() string { return p.Key }
func (p *HttpYamlParams) GetUnmarshaller() func(data []byte, v interface{}) error {
	return yaml.Unmarshal
}
