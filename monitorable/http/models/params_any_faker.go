//+build faker

package models

import "github.com/monitoror/monitoror/models"

type (
	HTTPAnyParams struct {
		URL           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`

		Status  models.TileStatus `json:"status" query:"status"`
		Message string            `json:"message" query:"message"`
		Values  []float64         `json:"values" query:"values"`
	}
)

func (p *HTTPAnyParams) IsValid() bool {
	return isValid(p.URL, p)
}

func (p *HTTPAnyParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPAnyParams) GetStatus() models.TileStatus { return p.Status }
func (p *HTTPAnyParams) GetMessage() string           { return p.Message }
func (p *HTTPAnyParams) GetValues() []float64         { return p.Values }
