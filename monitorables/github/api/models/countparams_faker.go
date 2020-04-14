//+build faker

package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	CountParams struct {
		Query string `json:"query" query:"query"`

		ValueValues []string `json:"valueValues" query:"valueValues"`
	}
)

func (p *CountParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.Query == "" {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}
