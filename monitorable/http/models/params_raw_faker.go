//+build faker

package models

import (
	"regexp"

	"github.com/monitoror/monitoror/models"
)

type (
	HttpRawParams struct {
		Url           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`
		Regex         string `json:"regex" query:"regex"`

		Status  models.TileStatus `json:"status" query:"status"`
		Message string            `json:"message" query:"message"`
		Values  []float64         `json:"values" query:"values"`
	}
)

func (p *HttpRawParams) IsValid() bool {
	if !isValid(p.Url, p) {
		return false
	}

	return isValidRegex(p)
}

func (p *HttpRawParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HttpRawParams) GetRegex() string          { return p.Regex }
func (p *HttpRawParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }

func (p *HttpRawParams) GetStatus() models.TileStatus { return p.Status }
func (p *HttpRawParams) GetMessage() string           { return p.Message }
func (p *HttpRawParams) GetValues() []float64         { return p.Values }
