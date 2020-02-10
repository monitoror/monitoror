//+build faker

package models

type (
	IssuesParams struct {
		Query string `json:"query" query:"query"`

		Values []float64 `json:"values" query:"values"`
	}
)

func (p *IssuesParams) IsValid() bool {
	return p.Query != ""
}
