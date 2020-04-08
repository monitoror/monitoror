//+build faker

package models

import (
	"fmt"
	"time"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	BuildParams struct {
		Job    string `json:"job" query:"job"`
		Branch string `json:"branch" query:"branch"`

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

func (p *BuildParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.Job == "" {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}

// Used by cache as identifier
func (p *BuildParams) String() string {
	return fmt.Sprintf("BUILD-%s-%s", p.Job, p.Branch)
}
