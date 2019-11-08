//+build faker

package models

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/models"
)

type (
	BuildParams struct {
		Job    string `json:"job" query:"job"`
		Branch string `json:"branch" query:"branch"`

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

func (p *BuildParams) IsValid() bool {
	return p.Job != ""
}

// Used by cache as identifier
func (p *BuildParams) String() string {
	return fmt.Sprintf("BUILD-%s-%s", p.Job, p.Branch)
}
