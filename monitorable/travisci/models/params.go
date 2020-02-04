//+build !faker

package models

import "fmt"

type (
	BuildParams struct {
		Owner      string `json:"owner" query:"owner"`
		Repository string `json:"repository" query:"repository"`
		Branch     string `json:"branch" query:"branch"`
	}
)

func (p *BuildParams) IsValid() bool {
	return p.Owner != "" && p.Repository != "" && p.Branch != ""
}

// Used by cache as identifier
func (p *BuildParams) String() string {
	return fmt.Sprintf("BUILD-%s-%s-%s", p.Owner, p.Repository, p.Branch)
}
