//+build faker

package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	CountParams struct {
		Query string `json:"query" query:"query"`

		ValueValues []string `json:"valueValues" query:"valueValues"`
	}
)

func (p *CountParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if p.Query == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "query" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "query"},
		}
	}

	return nil
}
