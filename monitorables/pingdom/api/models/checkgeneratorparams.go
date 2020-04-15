package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	CheckGeneratorParams struct {
		Tags   string `json:"tags,omitempty" query:"tags"`
		SortBy string `json:"sortBy,omitempty" query:"sortBy"`
	}
)

func (p *CheckGeneratorParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if p.SortBy != "" && p.SortBy != "name" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorInvalidFieldValue,
			Message: fmt.Sprintf(`Required "sortBy" field is missing.`),
			Data: uiConfigModels.ConfigErrorData{
				FieldName: "sortBy",
				Expected:  "name",
			},
		}
	}

	return nil
}
