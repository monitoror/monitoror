//+build !faker

package models

import "fmt"

type (
	BuildParams struct {
		Job    string `json:"job" query:"job"`
		Branch string `json:"branch" query:"branch"`
	}
)

func (p *BuildParams) IsValid() bool {
	return p.Job != ""
}

// Used by cache as identifier
func (p *BuildParams) String() string {
	return fmt.Sprintf("BUILD-%s-%s", p.Job, p.Branch)
}
