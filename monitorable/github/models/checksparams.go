package models

import "fmt"

type (
	ChecksParams struct {
		Owner      string `json:"owner" query:"owner"`
		Repository string `json:"repository" query:"repository"`
		Ref        string `json:"ref" query:"ref"`
	}
)

func (p *ChecksParams) IsValid() bool {
	return p.Owner != "" && p.Repository != "" && p.Ref != ""
}

// Used by cache as identifier
func (p *ChecksParams) String() string {
	return fmt.Sprintf("CHECKS-%s-%s-%s", p.Owner, p.Repository, p.Ref)
}
