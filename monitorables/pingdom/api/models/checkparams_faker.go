//+build faker

package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	CheckParams struct {
		ID *int `json:"id" query:"id"`

		Status coreModels.TileStatus `json:"status" query:"status"`
	}
)

func (p *CheckParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if p.ID == nil {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "id" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "id"},
		}
	}

	return nil
}
