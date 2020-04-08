//+build !faker

package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	PingParams struct {
		Hostname string `json:"hostname" query:"hostname"`
	}
)

func (p *PingParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.Hostname == "" {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}
