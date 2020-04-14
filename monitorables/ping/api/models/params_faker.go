//+build faker

package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	PingParams struct {
		Hostname string `json:"hostname" query:"hostname"`

		Status      coreModels.TileStatus `json:"status" query:"status"`
		ValueValues []string              `json:"valueValues" query:"valueValues"`
	}
)

func (p *PingParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.Hostname == "" {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}
