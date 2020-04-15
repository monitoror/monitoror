package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type PullRequestGeneratorParams struct {
	Owner      string `json:"owner" query:"owner"`
	Repository string `json:"repository" query:"repository"`
}

func (p *PullRequestGeneratorParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if p.Owner == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "owner" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "owner"},
		}
	}

	if p.Repository == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "repository" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "repository"},
		}
	}

	return nil
}
