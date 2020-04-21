package models

import "time"

type (
	TileBuild struct {
		PreviousStatus TileStatus `json:"previousStatus,omitempty"`

		ID           *string           `json:"id,omitempty"`
		Branch       *string           `json:"branch,omitempty"`
		MergeRequest *TileMergeRequest `json:"mergeRequest,omitempty"`
		Author       *Author           `json:"author,omitempty"`

		Duration          *int64     `json:"duration,omitempty"`          // In Seconds
		EstimatedDuration *int64     `json:"estimatedDuration,omitempty"` // In Seconds
		StartedAt         *time.Time `json:"startedAt,omitempty"`
		FinishedAt        *time.Time `json:"finishedAt,omitempty"`
	}

	TileMergeRequest struct {
		ID    int    `json:"id"`
		Title string `json:"title,omitempty"`
	}
)

func (t *Tile) WithBuild() *Tile {
	t.Build = &TileBuild{}
	return t
}
