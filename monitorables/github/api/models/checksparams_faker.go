//+build faker

package models

import (
	"fmt"
	"time"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	ChecksParams struct {
		Owner      string `json:"owner" query:"owner"`
		Repository string `json:"repository" query:"repository"`
		Ref        string `json:"ref" query:"ref"`

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

func (p *ChecksParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.Owner == "" {
		return &uiConfigModels.ConfigError{}
	}

	if p.Repository == "" {
		return &uiConfigModels.ConfigError{}
	}

	if p.Ref == "" {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}

// Used by cache as identifier
func (p *ChecksParams) String() string {
	return fmt.Sprintf("CHECKS-%s-%s-%s", p.Owner, p.Repository, p.Ref)
}
