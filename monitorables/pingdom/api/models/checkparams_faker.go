//+build faker

package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	CheckParams struct {
		ID *int `json:"id" query:"id"`

		Status coreModels.TileStatus `json:"status" query:"status"`
	}
)

func (p *CheckParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.ID == nil {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}
