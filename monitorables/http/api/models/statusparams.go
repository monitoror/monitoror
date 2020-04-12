//+build !faker

package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	HTTPStatusParams struct {
		URL           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin,omitempty" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax,omitempty" query:"statusCodeMax"`
	}
)

func (p *HTTPStatusParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if err := validateURL(p); err != nil {
		return err
	}

	if err := validateStatusCode(p); err != nil {
		return err
	}

	return nil
}

func (p *HTTPStatusParams) GetURL() (url string) { return p.URL }
func (p *HTTPStatusParams) GetStatusCodes() (min int, max int) {
	return getStatusCodesWithDefault(p.StatusCodeMin, p.StatusCodeMax)
}
