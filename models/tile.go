package models

import (
	"fmt"
	"strings"
)

type (
	Tile struct {
		Type   TileType   `json:"type"`
		Status TileStatus `json:"status"`

		Label   string `json:"label,omitempty"`
		Message string `json:"message,omitempty"`

		Value *TileValue `json:"value,omitempty"`
		Build *TileBuild `json:"build,omitempty"`
	}

	TileType   string //PING, PORT, ... (defined in usecase.go for each monitorable)
	TileStatus string
)

const (
	ActionRequiredStatus TileStatus = "ACTION_REQUIRED"
	CanceledStatus       TileStatus = "CANCELED"
	DisabledStatus       TileStatus = "DISABLED"
	FailedStatus         TileStatus = "FAILURE"
	QueuedStatus         TileStatus = "QUEUED"
	RunningStatus        TileStatus = "RUNNING"
	SuccessStatus        TileStatus = "SUCCESS"
	UnknownStatus        TileStatus = "UNKNOWN"
	WarningStatus        TileStatus = "WARNING"
)

const generatorPrefix string = "GENERATE:"

var AvailableTileStatuses = map[TileStatus]bool{
	ActionRequiredStatus: true,
	CanceledStatus:       true,
	DisabledStatus:       true,
	FailedStatus:         true,
	QueuedStatus:         true,
	RunningStatus:        true,
	SuccessStatus:        true,
	UnknownStatus:        true,
	WarningStatus:        true,
}

func NewTile(t TileType) *Tile {
	return &Tile{Type: t}
}

func NewGeneratorTileType(t TileType) TileType {
	tileType := fmt.Sprintf("%s%s", generatorPrefix, t)
	return TileType(tileType)
}

func (t TileType) IsGenerator() bool {
	return strings.HasPrefix(string(t), generatorPrefix)
}

func (t TileType) GetGeneratedTileType() TileType {
	return TileType(strings.TrimPrefix(string(t), generatorPrefix))
}
