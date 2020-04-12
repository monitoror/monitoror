package models

import (
	"fmt"
	"regexp"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	BuildGeneratorParams struct {
		Job string `json:"job" query:"job"`

		// Using Match / Unmatch filter instead of one filter because Golang's standard regex library doesn't have negative look ahead.
		Match   string `json:"match,omitempty" query:"match"`
		Unmatch string `json:"unmatch,omitempty" query:"unmatch"`
	}
)

func (p *BuildGeneratorParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if p.Job == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "job" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "job"},
		}
	}

	if p.Match != "" {
		if _, err := regexp.Compile(p.Match); err != nil {
			return &uiConfigModels.ConfigError{
				ID:      uiConfigModels.ConfigErrorInvalidFieldValue,
				Message: fmt.Sprintf(`Invalid "match" field. Must be a valid golang regex.`),
				Data: uiConfigModels.ConfigErrorData{
					FieldName: "match",
					Expected:  "valid golang regex",
				},
			}
		}
	}

	if p.Unmatch != "" {
		if _, err := regexp.Compile(p.Unmatch); err != nil {
			return &uiConfigModels.ConfigError{
				ID:      uiConfigModels.ConfigErrorInvalidFieldValue,
				Message: fmt.Sprintf(`Invalid "unmatch" field. Must be a valid golang regex.`),
				Data: uiConfigModels.ConfigErrorData{
					FieldName: "unmatch",
					Expected:  "valid golang regex",
				},
			}
		}
	}

	return nil
}
