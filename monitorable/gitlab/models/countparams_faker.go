//+build faker

package models

type (
	CountParams struct {
		Query string `json:"query" query:"query"`

		ValueValues []string `json:"valueValues" query:"valueValues"`
	}
)

func (p *CountParams) IsValid() bool {
	return p.Query != ""
}
