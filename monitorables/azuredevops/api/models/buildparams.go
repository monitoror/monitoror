//+build !faker

package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	BuildParams struct {
		Project    string  `json:"project" query:"project"`
		Definition *int    `json:"definition" query:"definition"`
		Branch     *string `json:"branch,omitempty" query:"branch"`
	}
)

func (p *BuildParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
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
func (p *BuildParams) String() string {
	str := fmt.Sprintf("BUILD-%s-%d", p.Project, *p.Definition)

	if p.Branch != nil {
		str = fmt.Sprintf("%s-%s", str, *p.Branch)
	}

	return str
}
