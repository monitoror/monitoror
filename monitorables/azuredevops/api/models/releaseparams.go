//+build !faker

package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	ReleaseParams struct {
		Project    string `json:"project" query:"project"`
		Definition *int   `json:"definition" query:"definition"`
	}
)

func (p *ReleaseParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if p.Project == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "project" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "project"},
		}
	}

	if p.Definition == nil {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "definition" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "definition"},
		}
	}

	return nil
}

// Used by cache as identifier
func (p *ReleaseParams) String() string {
	return fmt.Sprintf("RELEASE-%s-%d", p.Project, *p.Definition)
}
