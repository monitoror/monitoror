//+build !faker

package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	PortParams struct {
		Hostname string `json:"hostname" query:"hostname"`
		Port     int    `json:"port" query:"port"`
	}
)

func (p *PortParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	if p.Hostname == "" {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "hostname" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "hostname"},
		}
	}

	if p.Port == 0 {
		return &uiConfigModels.ConfigError{
			ID:      uiConfigModels.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "port" field is missing.`),
			Data:    uiConfigModels.ConfigErrorData{FieldName: "port"},
		}
	}

	return nil
}
