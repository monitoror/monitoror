//+build !faker

package models

type (
	JobParams struct {
		Job    string `json:"job" query:"job"`
		Parent string `json:"parent" query:"parent"`
	}
)

func (p *JobParams) IsValid() bool {
	return p.Job != ""
}
