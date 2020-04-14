package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	CheckGeneratorParams struct {
		Tags   string `json:"tags,omitempty" query:"tags"`
		SortBy string `json:"sortBy,omitempty" query:"sortBy"`
	}
)

func (p *CheckGeneratorParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.SortBy != "" && p.SortBy != "name" {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}
