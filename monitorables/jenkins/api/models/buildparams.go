//+build !faker

package models

import (
	"fmt"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	BuildParams struct {
		params.Default

		Job    string `json:"job" query:"job" validate:"required"`
		Branch string `json:"branch,omitempty" query:"branch"`
	}
)

// Used by cache as identifier
func (p *BuildParams) String() string {
	return fmt.Sprintf("BUILD-%s-%s", p.Job, p.Branch)
}
