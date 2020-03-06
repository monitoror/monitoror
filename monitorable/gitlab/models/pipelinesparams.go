//+build !faker

package models

import "fmt"

type (
	PipelinesParams struct {
		Repository string `json:"repository" query:"repository"`
		Ref        string `json:"ref" query:"ref"`
	}
)

func (p *PipelinesParams) IsValid() bool {
	return p.Repository != "" && p.Ref != ""
}

// Used by cache as identifier
func (p *PipelinesParams) String() string {
	return fmt.Sprintf("PIPELINES-%s-%s", p.Repository, p.Ref)
}
