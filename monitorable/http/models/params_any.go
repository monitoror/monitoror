package models

type (
	HttpAnyParams struct {
		Url           string `json:"url" query:"url"`
		StatusCodeMin *int   `json:"statusCodeMin" query:"statusCodeMin"`
		StatusCodeMax *int   `json:"statusCodeMax" query:"statusCodeMax"`
	}
)

func (p *HttpAnyParams) IsValid() bool {
	return isValid(p.Url, p)
}

func (p *HttpAnyParams) GetStatusCodes() (min int, max int) {
	return getStatusCodes(p.StatusCodeMin, p.StatusCodeMax)
}
