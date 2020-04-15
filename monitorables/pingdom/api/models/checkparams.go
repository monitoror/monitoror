//+build !faker

package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	CheckParams struct {
		ID *int `json:"id" query:"id"`
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
