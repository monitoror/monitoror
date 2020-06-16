//+build faker

package models

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	PipelineParams struct {
		params.Default

		ProjectID *int   `json:"projectId" query:"projectId" validate:"required"`
		Ref       string `json:"ref" query:"ref" validate:"required"`

		AuthorName      string `json:"authorName" query:"authorName"`
		AuthorAvatarURL string `json:"authorAvatarURL" query:"authorAvatarURL"`

		Status            coreModels.TileStatus `json:"status" query:"status"`
		PreviousStatus    coreModels.TileStatus `json:"previousStatus" query:"previousStatus"`
		StartedAt         time.Time             `json:"startedAt" query:"startedAt"`
		FinishedAt        time.Time             `json:"finishedAt" query:"finishedAt"`
		Duration          int64                 `json:"duration" query:"duration"`
		EstimatedDuration int64                 `json:"estimatedDuration" query:"estimatedDuration"`
	}
)

// Used by cache as identifier
func (p *PipelineParams) String() string {
	return fmt.Sprintf("PIPELINE-%d-%s", *p.ProjectID, p.Ref)
}
