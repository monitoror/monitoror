//+build faker

package models

import "github.com/monitoror/monitoror/models"

type (
	HttpAnyParams struct {
		Url           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`

		Status  models.TileStatus `json:"status" query:"status"`
		Message string            `json:"message" query:"message"`
		Values  []float64         `json:"values" query:"values"`
	}
)

func (p *HttpAnyParams) IsValid() bool {
	return isValid(p.Url, p)
}

func (p *HttpAnyParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HttpAnyParams) GetStatus() models.TileStatus { return p.Status }
func (p *HttpAnyParams) GetMessage() string           { return p.Message }
func (p *HttpAnyParams) GetValues() []float64         { return p.Values }
