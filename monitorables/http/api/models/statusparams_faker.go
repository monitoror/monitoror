//+build faker

package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	HTTPStatusParams struct {
		URL           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`

		Status  coreModels.TileStatus `json:"status" query:"status"`
		Message string                `json:"message" query:"message"`
	}
)

func (p *HTTPStatusParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if !isValid(p.URL, p) {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}

func (p *HTTPStatusParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPStatusParams) GetStatus() coreModels.TileStatus        { return p.Status }
func (p *HTTPStatusParams) GetMessage() string                      { return p.Message }
func (p *HTTPStatusParams) GetValueValues() []string                { panic("unimplemented") }
func (p *HTTPStatusParams) GetValueUnit() coreModels.TileValuesUnit { panic("unimplemented") }
