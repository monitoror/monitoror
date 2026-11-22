//go:build faker

package models

import (
	"fmt"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	CountParams struct {
		params.Default

		Owner      string `json:"owner" query:"owner" validate:"required"`
		Repository string `json:"repository" query:"repository" validate:"required"`
		Type       string `json:"type,omitempty" query:"type" validate:"omitempty,oneof=issues pulls"`
		State      string `json:"state,omitempty" query:"state" validate:"omitempty,oneof=open closed all"`
		Labels     string `json:"labels,omitempty" query:"labels"`
		Query      string `json:"query,omitempty" query:"query"`

		ValueValues []string `json:"valueValues" query:"valueValues"`
	}
)

// Used by cache as identifier
func (p *CountParams) String() string {
	return fmt.Sprintf("COUNT-%s-%s-%s-%s-%s-%s", p.Owner, p.Repository, p.Type, p.State, p.Labels, p.Query)
}
