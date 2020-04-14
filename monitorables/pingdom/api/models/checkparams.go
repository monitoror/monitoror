//+build !faker

package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	CheckParams struct {
		ID *int `json:"id" query:"id"`
	}
)

func (p *CheckParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.ID == nil {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}
