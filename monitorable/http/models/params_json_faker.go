//+build faker

package models

import (
	"encoding/json"
	"regexp"

	"github.com/monitoror/monitoror/models"
)

type (
	HTTPJsonParams struct {
		URL           string `json:"url" query:"url"`
		Key           string `json:"key" query:"key"`
		Regex         string `json:"regex" query:"regex"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`

		Status  models.TileStatus `json:"status" query:"status"`
		Message string            `json:"message" query:"message"`
		Values  []float64         `json:"values" query:"values"`
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

func (p *HTTPJsonParams) GetStatus() models.TileStatus { return p.Status }
func (p *HTTPJsonParams) GetMessage() string           { return p.Message }
func (p *HTTPJsonParams) GetValues() []float64         { return p.Values }
