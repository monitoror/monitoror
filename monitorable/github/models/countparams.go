//+build !faker

package models

type (
	CountParams struct {
		Query string `json:"query" query:"query"`
	}
)

func (p *CountParams) IsValid() bool {
	return p.Query != ""
}
