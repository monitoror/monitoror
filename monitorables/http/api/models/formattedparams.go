//+build !faker

package models

import (
	"regexp"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	HTTPFormattedParams struct {
		URL           string `json:"url" query:"url"`
		Format        Format `json:"format" query:"format"`
		Key           string `json:"key" query:"key"`
		Regex         string `json:"regex,omitempty" query:"regex"`
		StatusCodeMin *int   `json:"statusCodeMin,omitempty" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax,omitempty" query:"statusCodeMax"`
	}
)

func (p *HTTPFormattedParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if err := validateURL(p); err != nil {
		return err
	}

	if err := validateStatusCode(p); err != nil {
		return err
	}

	if err := validateFormat(p); err != nil {
		return err
	}

	if err := validateKey(p); err != nil {
		return err
	}

	if err := validateRegex(p); err != nil {
		return err
	}

	return nil
}

func (p *HTTPFormattedParams) GetURL() (url string) { return p.URL }
func (p *HTTPFormattedParams) GetStatusCodes() (min int, max int) {
	return getStatusCodesWithDefault(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPFormattedParams) GetRegex() string          { return p.Regex }
func (p *HTTPFormattedParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }

func (p *HTTPFormattedParams) GetKey() string    { return p.Key }
func (p *HTTPFormattedParams) GetFormat() Format { return p.Format }
