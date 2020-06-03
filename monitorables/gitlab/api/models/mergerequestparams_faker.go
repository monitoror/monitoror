//+build faker

package models

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	MergeRequestParams struct {
		params.Default

		ProjectID *int `json:"projectId" query:"projectId" validate:"required"`
		ID        *int `json:"id" query:"id" validate:"required"`

		Branch string `json:"branch" query:"branch"`

		AuthorName      string `json:"authorName" query:"authorName"`
		AuthorAvatarURL string `json:"authorAvatarURL" query:"authorAvatarURL"`

		MergeRequestTitle string `json:"mergeRequestTitle" query:"mergeRequestTitle"`

		Status            coreModels.TileStatus `json:"status" query:"status"`
		PreviousStatus    coreModels.TileStatus `json:"previousStatus" query:"previousStatus"`
		StartedAt         time.Time             `json:"startedAt" query:"startedAt"`
		FinishedAt        time.Time             `json:"finishedAt" query:"finishedAt"`
		Duration          int64                 `json:"duration" query:"duration"`
		EstimatedDuration int64                 `json:"estimatedDuration" query:"estimatedDuration"`
	}
)

// Used by cache as identifier
func (p *MergeRequestParams) String() string {
	return fmt.Sprintf("MERGEREQUEST-%d-%d", *p.ProjectID, *p.ID)
}
