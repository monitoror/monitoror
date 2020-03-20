//+build faker

package models

import coreModels "github.com/monitoror/monitoror/models"

type (
	HTTPStatusParams struct {
		URL           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`

		Status  coreModels.TileStatus `json:"status" query:"status"`
		Message string                `json:"message" query:"message"`
	}
)

func (p *HTTPStatusParams) IsValid() bool {
	return isValid(p.URL, p)
}

func (p *HTTPStatusParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPStatusParams) GetStatus() coreModels.TileStatus        { return p.Status }
func (p *HTTPStatusParams) GetMessage() string                      { return p.Message }
func (p *HTTPStatusParams) GetValueValues() []string                { panic("unimplemented") }
func (p *HTTPStatusParams) GetValueUnit() coreModels.TileValuesUnit { panic("unimplemented") }
