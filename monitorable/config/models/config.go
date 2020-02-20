package models

import "github.com/monitoror/monitoror/models"

type (
	Config struct {
		Version  *int     `json:"version"`
		Columns  *int     `json:"columns"`
		Tiles    []Tile   `json:"tiles,omitempty"`
		Errors   []string `json:"errors,omitempty"`
		Warnings []string `json:"warnings,omitempty"`
	}

	Tile struct {
		Type   models.TileType        `json:"type"`
		Params map[string]interface{} `json:"params,omitempty"`

		Label      string `json:"label,omitempty"`
		RowSpan    *int   `json:"rowSpan,omitempty"`
		ColumnSpan *int   `json:"columnSpan,omitempty"`

		Tiles           []Tile `json:"tiles,omitempty"`
		URL             string `json:"url,omitempty"`
		InitialMaxDelay *int   `json:"initialMaxDelay,omitempty"`

		// Used by config.hydrate only (will be remove before returning config to UI)
		ConfigVariant string `json:"configVariant,omitempty"`
	}
)

func (c *Config) AddErrors(reasons ...string) {
	c.Errors = append(c.Errors, reasons...)
}

func (c *Config) AddWarnings(reasons ...string) {
	c.Warnings = append(c.Warnings, reasons...)
}
