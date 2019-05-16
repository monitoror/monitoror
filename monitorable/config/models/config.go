package models

type (
	Config struct {
		Columns    int                      `json:"columns"`
		ApiBaseUrl string                   `json:"apiBaseUrl,omitempty"`
		Tiles      []map[string]interface{} `json:"tiles"`
	}
)
