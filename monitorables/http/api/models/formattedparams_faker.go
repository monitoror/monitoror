//+build faker

package models

import (
	"regexp"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	HTTPFormattedParams struct {
		URL           string `json:"url" query:"url"`
		Format        string `json:"format" query:"format"`
		Key           string `json:"key" query:"key"`
		Regex         string `json:"regex" query:"regex"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`

		Status      coreModels.TileStatus     `json:"status" query:"status"`
		Message     string                    `json:"message" query:"message"`
		ValueValues []string                  `json:"valueValues" query:"valueValues"`
		ValueUnit   coreModels.TileValuesUnit `json:"valueUnit" query:"valueUnit"`
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

func (p *HTTPFormattedParams) GetStatus() coreModels.TileStatus        { return p.Status }
func (p *HTTPFormattedParams) GetMessage() string                      { return p.Message }
func (p *HTTPFormattedParams) GetValueValues() []string                { return p.ValueValues }
func (p *HTTPFormattedParams) GetValueUnit() coreModels.TileValuesUnit { return p.ValueUnit }
