//+build !faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/validator"
)

type (
	HTTPStatusParams struct {
		URL           string `json:"url" query:"url" validate:"required,url,http"`
		StatusCodeMin *int   `json:"statusCodeMin,omitempty" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax,omitempty" query:"statusCodeMax"`
	}
)

func (p *HTTPStatusParams) Validate() []validator.Error {
	return validateStatusCode(p)
}

func (p *HTTPStatusParams) GetURL() (url string) { return p.URL }
func (p *HTTPStatusParams) GetStatusCodes() (min int, max int) {
	return getStatusCodesWithDefault(p.StatusCodeMin, p.StatusCodeMax)
}
