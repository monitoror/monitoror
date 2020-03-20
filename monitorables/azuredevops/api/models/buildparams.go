//+build !faker

package models

import "fmt"

type (
	BuildParams struct {
		Project    string  `json:"project" query:"project"`
		Definition *int    `json:"definition" query:"definition"`
		Branch     *string `json:"branch" query:"branch"`
	}
)

func (p *BuildParams) IsValid() bool {
	return p.Project != "" && p.Definition != nil
}

// Used by cache as identifier
func (p *BuildParams) String() string {
	str := fmt.Sprintf("BUILD-%s-%d", p.Project, *p.Definition)

	if p.Branch != nil {
		str = fmt.Sprintf("%s-%s", str, *p.Branch)
	}

	return str
}
