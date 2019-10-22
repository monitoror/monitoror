//+build !faker

package models

import "fmt"

type (
	BuildParams struct {
		Group      string `json:"group" query:"group"`
		Repository string `json:"repository" query:"repository"`
		Branch     string `json:"branch" query:"branch"`
	}
)

func (p *BuildParams) IsValid() bool {
	return p.Group != "" && p.Repository != "" && p.Branch != ""
}

// Used by cache as identifier
func (p *BuildParams) String() string {
	return fmt.Sprintf("BUILD-%s-%s-%s", p.Group, p.Repository, p.Branch)
}
