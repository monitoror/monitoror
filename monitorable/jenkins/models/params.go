//+build !faker

package models

type (
	BuildParams struct {
		Job    string `json:"job" query:"job"`
		Parent string `json:"parent" query:"parent"`
	}
)

func (p *BuildParams) IsValid() bool {
	return p.Job != ""
}
