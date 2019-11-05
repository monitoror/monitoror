//+build !faker

package models

import "fmt"

type (
	ReleaseParams struct {
		Project    string `json:"project" query:"project"`
		Definition *int   `json:"definition" query:"definition"`
	}
)

func (p *ReleaseParams) IsValid() bool {
	return p.Project != "" && p.Definition != nil
}

// Used by cache as identifier
func (p *ReleaseParams) String() string {
	return fmt.Sprintf("RELEASE-%s-%d", p.Project, *p.Definition)
}
