package models

import "time"

type (
	Tile struct {
		Type   TileType   `json:"type"`
		Status TileStatus `json:"status"`

		Label   string        `json:"label,omitempty"`
		Message string        `json:"message,omitempty"`
		Values  []float64     `json:"values,omitempty"`
		Unit    TileValueUnit `json:"unit,omitempty"`

		Author *Author `json:"author,omitempty"`

		PreviousStatus    TileStatus `json:"previousStatus,omitempty"`
		StartedAt         *time.Time `json:"startedAt,omitempty"`
		FinishedAt        *time.Time `json:"finishedAt,omitempty"`
		Duration          *int64     `json:"duration,omitempty"`          // In Seconds
		EstimatedDuration *int64     `json:"estimatedDuration,omitempty"` // In Seconds
	}

	TileType      string // PING, JENKINS_BUILD, ...
	TileStatus    string // SUCCESS, FAILURE, ...
	TileValueUnit string // MILLISECOND, NONE, ...

	Author struct {
		Name      string `json:"name,omitempty"`
		AvatarURL string `json:"avatarURL,omitempty"`
	}
)

// List of all Status Code
const (
	SuccessStatus  TileStatus = "SUCCESS"
	FailedStatus   TileStatus = "FAILURE"
	RunningStatus  TileStatus = "RUNNING"
	QueuedStatus   TileStatus = "QUEUED"
	DisabledStatus TileStatus = "DISABLED"
	CanceledStatus TileStatus = "CANCELED"
	WarningStatus  TileStatus = "WARNING"
	UnknownStatus  TileStatus = "UNKNOWN"
)

const (
	MillisecondUnit TileValueUnit = "MILLISECOND"
	DefaultUnit     TileValueUnit = ""
)

func NewTile(t TileType) *Tile {
	return &Tile{
		Type: t,
	}
}
