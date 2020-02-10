//+build faker

package models

type (
	CountParams struct {
		Query string `json:"query" query:"query"`

		Values []float64 `json:"values" query:"values"`
	}
)

func (p *CountParams) IsValid() bool {
	return p.Query != ""
}
