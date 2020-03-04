package models

import "github.com/monitoror/monitoror/models"

type (
	ConfigBag struct {
		Config *Config       `json:"config,omitempty"`
		Errors []ConfigError `json:"errors,omitempty"`
	}

	Config struct {
		Version *ConfigVersion `json:"version"`
		Columns *int           `json:"columns"`
		Zoom    *float32       `json:"zoom,omitempty"`
		Tiles   []Tile         `json:"tiles"`
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

		// Used by config.hydrate only (will be removed before returning config to UI)
		ConfigVariant string `json:"configVariant,omitempty"`
	}

	ConfigError struct {
		ID      ConfigErrorID   `json:"id"`
		Message string          `json:"message"`
		Data    ConfigErrorData `json:"data"`
	}

	ConfigErrorData struct {
		ConfigExtract          string `json:"configExtract,omitempty"`
		ConfigExtractHighlight string `json:"configExtractHighlight,omitempty"`
		Value                  string `json:"value,omitempty"`
		FieldName              string `json:"fieldName,omitempty"`
		Expected               string `json:"expected,omitempty"`
	}

	ConfigErrorID string
)

const (
	ConfigErrorConfigNotFound          ConfigErrorID = "ERROR_CONFIG_NOT_FOUND"
	ConfigErrorInvalidFieldValue       ConfigErrorID = "ERROR_INVALID_FIELD_VALUE"
	ConfigErrorMissingRequiredField    ConfigErrorID = "ERROR_MISSING_REQUIRED_FIELD"
	ConfigErrorUnauthorizedField       ConfigErrorID = "ERROR_UNAUTHORIZED_FIELD"
	ConfigErrorUnauthorizedSubtileType ConfigErrorID = "ERROR_UNAUTHORIZED_SUBTILE_TYPE"
	ConfigErrorUnableToHydrate         ConfigErrorID = "ERROR_UNABLE_TO_HYDRATE"
	ConfigErrorUnableToParseConfig     ConfigErrorID = "ERROR_UNABLE_TO_PARSE_CONFIG"
	ConfigErrorUnexpectedError         ConfigErrorID = "ERROR_UNEXPECTED"
	ConfigErrorUnknownTileType         ConfigErrorID = "ERROR_UNKNOWN_TILE_TYPE"
	ConfigErrorUnknownVariant          ConfigErrorID = "ERROR_UNKNOWN_VARIANT"
	ConfigErrorUnsupportedVersion      ConfigErrorID = "ERROR_UNSUPPORTED_VERSION"
)

func (c *ConfigBag) AddErrors(errors ...ConfigError) {
	c.Errors = append(c.Errors, errors...)
}
