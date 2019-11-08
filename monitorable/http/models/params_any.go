//+build !faker

package models

type (
	HTTPAnyParams struct {
		URL           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`
	}
)

func (p *HTTPAnyParams) IsValid() bool {
	return isValid(p.URL, p)
}

func (p *HTTPAnyParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}
