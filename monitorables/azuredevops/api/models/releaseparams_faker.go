//+build faker

package models

import (
	"time"

	"github.com/monitoror/monitoror/models"
)

type (
	ReleaseParams struct {
		Project    string `json:"project" query:"project"`
		Definition *int   `json:"definition" query:"definition"`

		AuthorName      string `json:"authorName" query:"authorName"`
		AuthorAvatarURL string `json:"authorAvatarURL" query:"authorAvatarURL"`

		Status            models.TileStatus `json:"status" query:"status"`
		PreviousStatus    models.TileStatus `json:"previousStatus" query:"previousStatus"`
		StartedAt         time.Time         `json:"startedAt" query:"startedAt"`
		FinishedAt        time.Time         `json:"finishedAt" query:"finishedAt"`
		Duration          int64             `json:"duration" query:"duration"`
		EstimatedDuration int64             `json:"estimatedDuration" query:"estimatedDuration"`
	}
)

func (p *ReleaseParams) IsValid() bool {
	return p.Project != "" && p.Definition != nil
}
