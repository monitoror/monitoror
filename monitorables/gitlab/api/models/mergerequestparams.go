//+build !faker

package models

import (
	"fmt"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	MergeRequestParams struct {
		params.Default

		ProjectID *int `json:"projectId" query:"projectId" validate:"required"`
		ID        *int `json:"id" query:"id" validate:"required"`
	}
)

// Used by cache as identifier
func (p *MergeRequestParams) String() string {
	return fmt.Sprintf("MERGEREQUEST-%d-%d", *p.ProjectID, *p.ID)
}
