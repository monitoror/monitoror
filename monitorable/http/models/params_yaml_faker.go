//+build faker

package models

import (
	"regexp"

	"github.com/monitoror/monitoror/models"
	"gopkg.in/yaml.v2"
)

type (
	HTTPYamlParams struct {
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

func (p *HTTPYamlParams) IsValid() bool {
	if !isValid(p.URL, p) {
		return false
	}

	if !isValidKey(p) {
		return false
	}

	return isValidRegex(p)
}

func (p *HTTPYamlParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPYamlParams) GetRegex() string          { return p.Regex }
func (p *HTTPYamlParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }

func (p *HTTPYamlParams) GetKey() string { return p.Key }
func (p *HTTPYamlParams) GetUnmarshaller() func(data []byte, v interface{}) error {
	return yaml.Unmarshal
}

func (p *HTTPYamlParams) GetStatus() models.TileStatus { return p.Status }
func (p *HTTPYamlParams) GetMessage() string           { return p.Message }
func (p *HTTPYamlParams) GetValues() []float64         { return p.Values }
