//+build !faker

package models

import (
	"fmt"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	ReleaseParams struct {
		params.Default

		Project    string `json:"project" query:"project" validate:"required"`
		Definition *int   `json:"definition" query:"definition" validate:"required"`
	}
)

// Used by cache as identifier
func (p *ReleaseParams) String() string {
	return fmt.Sprintf("RELEASE-%s-%d", p.Project, *p.Definition)
}
