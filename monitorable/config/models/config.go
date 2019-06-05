package models

type (
	Config struct {
		Columns int                      `json:"columns"`
		Tiles   []map[string]interface{} `json:"tiles"`
	}
)
