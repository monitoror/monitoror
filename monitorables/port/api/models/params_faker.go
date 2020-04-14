//+build faker

package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	PortParams struct {
		Hostname string `json:"hostname" query:"hostname"`
		Port     int    `json:"port" query:"port"`

		Status coreModels.TileStatus `json:"status" query:"status"`
	}
)

func (p *PortParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.Hostname == "" {
		return &uiConfigModels.ConfigError{}
	}

	if p.Port == 0 {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}
