//+build faker

package models

import (
	"regexp"

	"github.com/monitoror/monitoror/models"
)

type (
	HTTPRawParams struct {
		URL           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`
		Regex         string `json:"regex" query:"regex"`

		Status      models.TileStatus `json:"status" query:"status"`
		Message     string            `json:"message" query:"message"`
		ValueValues []string          `json:"valueValues" query:"valueValues"`
	}
)

func (p *HTTPRawParams) IsValid() bool {
	if !isValid(p.URL, p) {
		return false
	}

	return isValidRegex(p)
}

func (p *HTTPRawParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPRawParams) GetRegex() string          { return p.Regex }
func (p *HTTPRawParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }

func (p *HTTPRawParams) GetStatus() models.TileStatus { return p.Status }
func (p *HTTPRawParams) GetMessage() string           { return p.Message }
func (p *HTTPRawParams) GetValueValues() []string     { return p.ValueValues }
