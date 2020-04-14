//+build faker

package models

import (
	"regexp"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	HTTPRawParams struct {
		URL           string `json:"url" query:"url"`
		Regex         string `json:"regex,omitempty" query:"regex"`
		StatusCodeMin *int   `json:"statusCodeMin,omitempty" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax,omitempty" query:"statusCodeMax"`

		Status      coreModels.TileStatus     `json:"status" query:"status"`
		Message     string                    `json:"message" query:"message"`
		ValueValues []string                  `json:"valueValues" query:"valueValues"`
		ValueUnit   coreModels.TileValuesUnit `json:"valueUnit" query:"valueUnit"`
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

func (p *HTTPRawParams) GetStatus() coreModels.TileStatus        { return p.Status }
func (p *HTTPRawParams) GetMessage() string                      { return p.Message }
func (p *HTTPRawParams) GetValueValues() []string                { return p.ValueValues }
func (p *HTTPRawParams) GetValueUnit() coreModels.TileValuesUnit { return p.ValueUnit }
