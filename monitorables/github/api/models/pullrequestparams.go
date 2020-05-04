//+build !faker

package models

import (
	"fmt"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	PullRequestParams struct {
		params.Default

		Owner      string `json:"owner" query:"owner" validate:"required"`
		Repository string `json:"repository" query:"repository" validate:"required"`
		ID         *int   `json:"id" query:"id" validate:"required"`
	}
)

// Used by cache as identifier
func (p *PullRequestParams) String() string {
	return fmt.Sprintf("PULLREQUEST-%s-%s-%d", p.Owner, p.Repository, *p.ID)
}
