//+build !faker

package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	PingParams struct {
		Hostname string `json:"hostname" query:"hostname"`
	}
)

func (p *PingParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if p.Hostname == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "hostname" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "hostname"},
		}
	}

	return nil
}
