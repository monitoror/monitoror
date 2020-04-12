//+build !faker

package models

import (
	"regexp"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	HTTPRawParams struct {
		URL           string `json:"url" query:"url"`
		Regex         string `json:"regex,omitempty" query:"regex"`
		StatusCodeMin *int   `json:"statusCodeMin,omitempty" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax,omitempty" query:"statusCodeMax"`
	}
)

func (p *HTTPRawParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if err := validateURL(p); err != nil {
		return err
	}

	if err := validateStatusCode(p); err != nil {
		return err
	}

	if err := validateRegex(p); err != nil {
		return err
	}

	return nil
}

func (p *HTTPRawParams) GetURL() (url string) { return p.URL }
func (p *HTTPRawParams) GetStatusCodes() (min int, max int) {
	return getStatusCodesWithDefault(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPRawParams) GetRegex() string          { return p.Regex }
func (p *HTTPRawParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }
