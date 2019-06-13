package models

type (
	Config struct {
		Version int                      `json:"version"`
		Columns int                      `json:"columns"`
		Tiles   []map[string]interface{} `json:"tiles"`
	}
)
