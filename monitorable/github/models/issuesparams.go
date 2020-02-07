package models

type (
	IssuesParams struct {
		Query string `json:"query" query:"query"`
	}
)

func (p *IssuesParams) IsValid() bool {
	return p.Query != ""
}
