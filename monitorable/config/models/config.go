package models

import (
	"github.com/monitoror/monitoror/models/tiles"
)

type (
	Config struct {
		Version int    `json:"version"`
		Columns int    `json:"columns"`
		Tiles   []Tile `json:"tiles"`
	}

	Tile struct {
		Type   tiles.TileType         `json:"type"`
		Params map[string]interface{} `json:"params,omitempty"`

		Label      string `json:"label,omitempty"`
		RowSpan    *int   `json:"rowSpan,omitempty"`
		ColumnSpan *int   `json:"columnSpan,omitempty"`

		Tiles []Tile `json:"tiles,omitempty"`
		Url   string `json:"url,omitempty"`

		// Used by config.hydrate only (will be remove before returning config to front)
		ConfigVariant string `json:"configVariant,omitempty"`
	}
)
