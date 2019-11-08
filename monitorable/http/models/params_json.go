//+build !faker

package models

import (
	"encoding/json"
	"regexp"
)

type (
	HTTPJsonParams struct {
		URL           string `json:"url" query:"url"`
		Key           string `json:"key" query:"key"`
		Regex         string `json:"regex" query:"regex"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`
	}
)

func (p *HTTPJsonParams) IsValid() bool {
	if !isValid(p.URL, p) {
		return false
	}

	if !isValidKey(p) {
		return false
	}

	return isValidRegex(p)
}

func (p *HTTPJsonParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPJsonParams) GetRegex() string          { return p.Regex }
func (p *HTTPJsonParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }

func (p *HTTPJsonParams) GetKey() string { return p.Key }
func (p *HTTPJsonParams) GetUnmarshaller() func(data []byte, v interface{}) error {
	return json.Unmarshal
}
