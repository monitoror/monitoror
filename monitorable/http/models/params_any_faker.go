//+build faker

package models

import "github.com/monitoror/monitoror/models/tiles"

type (
	HttpAnyParams struct {
		Url           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`

		Status  tiles.TileStatus `json:"status" query:"status"`
		Message string           `json:"message" query:"message"`
	}
)

func (p *HttpAnyParams) IsValid() bool {
	return isValid(p.Url, p)
}

func (p *HttpAnyParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HttpAnyParams) GetStatus() tiles.TileStatus { return p.Status }
func (p *HttpAnyParams) GetMessage() string          { return p.Message }
