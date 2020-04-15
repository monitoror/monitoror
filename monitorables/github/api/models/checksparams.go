//+build !faker

package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	ChecksParams struct {
		Owner      string `json:"owner" query:"owner"`
		Repository string `json:"repository" query:"repository"`
		Ref        string `json:"ref" query:"ref"`
	}
)

func (p *ChecksParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if p.Owner == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "repository" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "repository"},
		}
	}

	if p.Repository == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "owner" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "owner"},
		}
	}

	if p.Ref == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "ref" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "ref"},
		}
	}

	return nil
}

// Used by cache as identifier
func (p *ChecksParams) String() string {
	return fmt.Sprintf("CHECKS-%s-%s-%s", p.Owner, p.Repository, p.Ref)
}
