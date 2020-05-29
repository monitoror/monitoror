//+build !faker

package models

import (
	"fmt"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	PipelineParams struct {
		params.Default

		ProjectID *int   `json:"projectId" query:"projectId" validate:"required"`
		Ref       string `json:"ref" query:"ref" validate:"required"`
	}
)

// Used by cache as identifier
func (p *PipelineParams) String() string {
	return fmt.Sprintf("PIPELINE-%d-%s", *p.ProjectID, p.Ref)
}
