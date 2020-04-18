//+build !faker

package models

import (
	"fmt"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	BuildParams struct {
		params.Default

		Project    string  `json:"project" query:"project" validate:"required"`
		Definition *int    `json:"definition" query:"definition" validate:"required"`
		Branch     *string `json:"branch,omitempty" query:"branch"`
	}
)

// Used by cache as identifier
func (p *BuildParams) String() string {
	str := fmt.Sprintf("BUILD-%s-%d", p.Project, *p.Definition)

	if p.Branch != nil {
		str = fmt.Sprintf("%s-%s", str, *p.Branch)
	}

	return str
}
