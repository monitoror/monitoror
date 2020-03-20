//+build !faker

package models

import "regexp"

type (
	HTTPRawParams struct {
		URL           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`
		Regex         string `json:"regex" query:"regex"`
	}
)

func (p *HTTPRawParams) IsValid() bool {
	if !isValid(p.URL, p) {
		return false
	}

	return isValidRegex(p)
}

func (p *HTTPRawParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HTTPRawParams) GetRegex() string          { return p.Regex }
func (p *HTTPRawParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }
