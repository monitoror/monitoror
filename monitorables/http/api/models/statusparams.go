//+build !faker

package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	HTTPStatusParams struct {
		URL           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin,omitempty" query:"statusCodeMin,omitempty"`
		StatusCodeMax *int   `json:"statusCodeMax,omitempty" query:"statusCodeMax,omitempty"`
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
