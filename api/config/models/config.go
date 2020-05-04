package models

import (
	"github.com/monitoror/monitoror/api/config/versions"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	ConfigBag struct {
		Config *Config       `json:"config,omitempty"`
		Errors []ConfigError `json:"errors,omitempty"`
	}

	Config struct {
		Version *versions.ConfigVersion `json:"version"`
		Columns *int                    `json:"columns" validate:"required,gt=0"`
		Zoom    *float32                `json:"zoom,omitempty" validate:"omitempty,gt=0,lte=10"`
		Tiles   []TileConfig            `json:"tiles" validate:"required,notempty"`
	}

	TileConfig struct {
		Type coreModels.TileType `json:"type" validate:"required"`

		Label      string `json:"label,omitempty"`
		RowSpan    *int   `json:"rowSpan,omitempty" validate:"omitempty,gt=0"`
		ColumnSpan *int   `json:"columnSpan,omitempty" validate:"omitempty,gt=0"`

		Tiles           []TileConfig `json:"tiles,omitempty"`
		URL             string       `json:"url,omitempty"`
		InitialMaxDelay *int         `json:"initialMaxDelay,omitempty"`

		// Used to validate config and to create API URLs
		// Will be removed before being returned to the UI
		Params        map[string]interface{} `json:"params,omitempty"`
		ConfigVariant coreModels.VariantName `json:"configVariant,omitempty"`
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
	ConfigErrorConfigNotFound                    ConfigErrorID = "ERROR_CONFIG_NOT_FOUND"
	ConfigErrorDisabledVariant                   ConfigErrorID = "ERROR_DISABLED_VARIANT"
	ConfigErrorFieldTypeMismatch                 ConfigErrorID = "ERROR_FIELD_TYPE_MISMATCH"
	ConfigErrorInvalidEscapedCharacter           ConfigErrorID = "ERROR_INVALID_ESCAPED_CHARACTER"
	ConfigErrorInvalidFieldValue                 ConfigErrorID = "ERROR_INVALID_FIELD_VALUE"
	ConfigErrorMissingRequiredField              ConfigErrorID = "ERROR_MISSING_REQUIRED_FIELD"
	ConfigErrorUnsupportedTileInThisVersion      ConfigErrorID = "ERROR_UNSUPPORTED_TILE_IN_THIS_VERSION"
	ConfigErrorUnsupportedTileParamInThisVersion ConfigErrorID = "ERROR_UNSUPPORTED_TILE_PARAM_IN_THIS_VERSION"
	ConfigErrorUnauthorizedField                 ConfigErrorID = "ERROR_UNAUTHORIZED_FIELD"
	ConfigErrorUnauthorizedSubtileType           ConfigErrorID = "ERROR_UNAUTHORIZED_SUBTILE_TYPE"
	ConfigErrorUnableToHydrate                   ConfigErrorID = "ERROR_UNABLE_TO_HYDRATE"
	ConfigErrorUnableToParseConfig               ConfigErrorID = "ERROR_UNABLE_TO_PARSE_CONFIG"
	ConfigErrorUnexpectedError                   ConfigErrorID = "ERROR_UNEXPECTED"
	ConfigErrorUnknownField                      ConfigErrorID = "ERROR_UNKNOWN_FIELD"
	ConfigErrorUnknownGeneratorTileType          ConfigErrorID = "ERROR_UNKNOWN_GENERATOR_TILE_TYPE"
	ConfigErrorUnknownNamedConfig                ConfigErrorID = "ERROR_UNKNOWN_NAMED_CONFIG"
	ConfigErrorUnknownTileType                   ConfigErrorID = "ERROR_UNKNOWN_TILE_TYPE"
	ConfigErrorUnknownVariant                    ConfigErrorID = "ERROR_UNKNOWN_VARIANT"
	ConfigErrorUnsupportedVersion                ConfigErrorID = "ERROR_UNSUPPORTED_VERSION"
)

func (c *ConfigBag) AddErrors(errors ...ConfigError) {
	c.Errors = append(c.Errors, errors...)
}
