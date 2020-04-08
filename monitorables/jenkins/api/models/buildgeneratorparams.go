package models

import (
	"regexp"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	BuildGeneratorParams struct {
		Job string `json:"job" query:"job"`

		// Using Match / Unmatch filter instead of one filter because Golang's standard regex library doesn't have negative look ahead.
		Match   string `json:"match,omitempty" query:"match,omitempty"`
		Unmatch string `json:"unmatch,omitempty" query:"unmatch,omitempty"`
	}
)

func (p *BuildGeneratorParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.Job == "" {
		return &uiConfigModels.ConfigError{}
	}

	if p.Match != "" {
		if _, err := regexp.Compile(p.Match); err != nil {
			return &uiConfigModels.ConfigError{}
		}
	}

	if p.Unmatch != "" {
		if _, err := regexp.Compile(p.Unmatch); err != nil {
			return &uiConfigModels.ConfigError{}
		}
	}

	return nil
}
