//+build !faker

package models

import (
	"regexp"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	HTTPRawParams struct {
		URL           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin,omitempty" query:"statusCodeMin,omitempty"`
		StatusCodeMax *int   `json:"statusCodeMax,omitempty" query:"statusCodeMax,omitempty"`
		Regex         string `json:"regex,omitempty" query:"regex,omitempty"`
	}
)

func (p *HTTPRawParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if !isValid(p.URL, p) {
		return &uiConfigModels.ConfigError{}
	}

	if !isValidRegex(p) {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}

func (p *HTTPRawParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPRawParams) GetRegex() string          { return p.Regex }
func (p *HTTPRawParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }
