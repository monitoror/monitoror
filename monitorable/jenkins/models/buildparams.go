//+build !faker

package models

type (
	BuildParams struct {
		Job    string `json:"job" query:"job"`
		Branch string `json:"branch" query:"branch"`
	}
)

func (p *BuildParams) IsValid() bool {
	return p.Job != ""
}
