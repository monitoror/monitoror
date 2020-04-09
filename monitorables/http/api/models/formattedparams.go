//+build !faker

package models

import (
	"regexp"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	HTTPFormattedParams struct {
		URL           string `json:"url" query:"url"`
		Format        string `json:"format" query:"format"`
		Key           string `json:"key" query:"key"`
		Regex         string `json:"regex,omitempty" query:"regex"`
		StatusCodeMin *int   `json:"statusCodeMin,omitempty" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax,omitempty" query:"statusCodeMax"`
	}
)

func (p *HTTPFormattedParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if !isValid(p.URL, p) {
		return &uiConfigModels.ConfigError{}
	}

	if !isSupportedFormat(p) {
		return &uiConfigModels.ConfigError{}
	}

	if !isValidKey(p) {
		return &uiConfigModels.ConfigError{}
	}

	if !isValidRegex(p) {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}

func (p *HTTPFormattedParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPFormattedParams) GetRegex() string          { return p.Regex }
func (p *HTTPFormattedParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }

func (p *HTTPFormattedParams) GetKey() string    { return p.Key }
func (p *HTTPFormattedParams) GetFormat() string { return p.Format }
