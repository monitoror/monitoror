package models

import "regexp"

type (
	HttpRawParams struct {
		Url           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`
		Regex         string `json:"regex" query:"regex"`
	}
)

func (p *HttpRawParams) IsValid() bool {
	if !isValid(p.Url, p) {
		return false
	}

	return isValidRegex(p)
}

func (p *HttpRawParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HttpRawParams) GetRegex() string          { return p.Regex }
func (p *HttpRawParams) GetRegexp() *regexp.Regexp { return getRegexp(p.GetRegex()) }
