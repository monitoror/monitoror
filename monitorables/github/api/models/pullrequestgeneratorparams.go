package models

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type PullRequestGeneratorParams struct {
	Owner      string `json:"owner" query:"owner"`
	Repository string `json:"repository" query:"repository"`
}

func (p *PullRequestGeneratorParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.Owner == "" {
		return &uiConfigModels.ConfigError{}
	}

	if p.Repository == "" {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}
