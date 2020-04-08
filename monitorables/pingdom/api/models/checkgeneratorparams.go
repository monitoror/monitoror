package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	CheckGeneratorParams struct {
		Tags   string `json:"tags,omitempty" query:"tags,omitempty"`
		SortBy string `json:"sortBy,omitempty" query:"sortBy,omitempty"`
	}
)

func (p *CheckGeneratorParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.SortBy != "" && p.SortBy != "name" {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}
