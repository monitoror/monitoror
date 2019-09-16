//+build !faker

package models

import "regexp"

type (
	// TODO : IMPROVEMENT: inheritance impossible due to echo context.Bind. Code custom binder for use struct inheritance
	HttpAnyParams struct {
		Url           string `json:"Url" query:"Url"`
		StatusCodeMin *int   `json:"StatusCodeMin" query:"StatusCodeMin"`
		StatusCodeMax *int   `json:"StatusCodeMax" query:"StatusCodeMax"`
	}

	HttpRawParams struct {
		Url           string `json:"Url" query:"Url"`
		Regex         string `json:"Regex" query:"Regex"`
		StatusCodeMin *int   `json:"StatusCodeMin" query:"StatusCodeMin"`
		StatusCodeMax *int   `json:"StatusCodeMax" query:"StatusCodeMax"`
	}

	// HttpFormattedDataParams : JSON / YAML
	HttpFormattedDataParams struct {
		Url           string `json:"Url" query:"Url"`
		Key           string `json:"Key" query:"Key"`
		Regex         string `json:"Regex" query:"Regex"`
		StatusCodeMin *int   `json:"StatusCodeMin" query:"StatusCodeMin"`
		StatusCodeMax *int   `json:"StatusCodeMax" query:"StatusCodeMax"`
	}

	StatusCodeRangeProvider interface {
		GetStatusCodeRange() (int, int)
	}

	RegexProvider interface {
		GetRegex() *regexp.Regexp
	}

	KeyProvider interface {
		GetKey() string
	}
)

const (
	DefaultMinStatusCode = 200
	DefaultMaxStatusCode = 399
)

func (p *HttpAnyParams) IsValid() bool {
	min, max := p.GetStatusCodeRange()
	return isValid(p.Url, min, max, "")
}
func (p *HttpAnyParams) GetStatusCodeRange() (min int, max int) {
	return getStatusCodeRange(p.StatusCodeMin, p.StatusCodeMax)
}

func (p *HttpRawParams) IsValid() bool {
	min, max := p.GetStatusCodeRange()
	return isValid(p.Url, min, max, p.Regex)
}
func (p *HttpRawParams) GetStatusCodeRange() (min int, max int) {
	return getStatusCodeRange(p.StatusCodeMin, p.StatusCodeMax)
}
func (p *HttpRawParams) GetRegex() *regexp.Regexp {
	if p.Regex != "" {
		return regexp.MustCompile(p.Regex)
	}
	return nil
}

func (p *HttpFormattedDataParams) IsValid() bool {
	key := p.GetKey()
	if key == "" || key == "." {
		return false
	}

	min, max := p.GetStatusCodeRange()
	return isValid(p.Url, min, max, p.Regex)
}
func (p *HttpFormattedDataParams) GetStatusCodeRange() (min int, max int) {
	return getStatusCodeRange(p.StatusCodeMin, p.StatusCodeMax)
}
func (p *HttpFormattedDataParams) GetRegex() *regexp.Regexp {
	if p.Regex != "" {
		return regexp.MustCompile(p.Regex)
	}
	return nil
}
func (p *HttpFormattedDataParams) GetKey() string {
	return p.Key
}

func isValid(url string, min, max int, regex string) bool {
	if url == "" {
		return false
	}

	if min > max {
		return false
	}

	if regex != "" {
		_, err := regexp.Compile(regex)
		if err != nil {
			return false
		}
	}

	return true
}

func getStatusCodeRange(statusCodeMin, statusCodeMax *int) (min int, max int) {
	min = DefaultMinStatusCode
	if statusCodeMin != nil {
		min = *statusCodeMin
	}
	max = DefaultMaxStatusCode
	if statusCodeMax != nil {
		max = *statusCodeMax
	}
	return
}
