package models

import "time"

type (
	Tile struct {
		Type   TileType   `json:"type"`
		Status TileStatus `json:"status"`

		Label   string        `json:"label,omitempty"`
		Message string        `json:"message,omitempty"`
		Values  []float64     `json:"value,omitempty"`
		Unit    TileValueUnit `json:"unit,omitempty"`

		Author *Author `json:"author,omitempty"`

		PreviousStatus    TileStatus `json:"previousStatus,omitempty"`
		StartedAt         *time.Time `json:"startedAt,omitempty"`
		FinishedAt        *time.Time `json:"finishedAt,omitempty"`
		Duration          *int64     `json:"duration,omitempty"`
		EstimatedDuration *int64     `json:"estimatedDuration,omitempty"`
	}

	TileType      string // PING, JENKINS_BUILD, ...
	TileStatus    string // SUCCESS, FAILURE, ...
	TileValueUnit string // MILLISECOND, NONE, ...

	Author struct {
		Name      string `json:"name,omitempty"`
		AvatarUrl string `json:"avatarUrl,omitempty"`
	}
)

// List of all Status Code
const (
	SuccessStatus  TileStatus = "SUCCESS"
	FailedStatus   TileStatus = "FAILURE"
	RunningStatus  TileStatus = "RUNNING"
	QueuedStatus   TileStatus = "QUEUED"
	DisabledStatus TileStatus = "DISABLED"
	AbortedStatus  TileStatus = "ABORTED"
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
