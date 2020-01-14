//+build faker

package models

import (
	"regexp"

	"github.com/monitoror/monitoror/models"
)

type (
	HTTPFormattedParams struct {
		URL           string `json:"url" query:"url"`
		Format        string `json:"format" query:"format"`
		Key           string `json:"key" query:"key"`
		Regex         string `json:"regex" query:"regex"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`

		Status  models.TileStatus `json:"status" query:"status"`
		Message string            `json:"message" query:"message"`
		Values  []float64         `json:"values" query:"values"`
	}
)

func (p *HTTPFormattedParams) IsValid() bool {
	if !isValid(p.URL, p) {
		return false
	}

	if !isSupportedFormat(p) {
		return false
	}

	if !isValidKey(p) {
		return false
	}

	return isValidRegex(p)
}

func (p *HTTPFormattedParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPFormattedParams) GetRegex() string          { return p.Regex }
func (p *HTTPFormattedParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }

func (p *HTTPFormattedParams) GetKey() string    { return p.Key }
func (p *HTTPFormattedParams) GetFormat() string { return p.Format }

func (p *HTTPFormattedParams) GetStatus() models.TileStatus { return p.Status }
func (p *HTTPFormattedParams) GetMessage() string           { return p.Message }
func (p *HTTPFormattedParams) GetValues() []float64         { return p.Values }
