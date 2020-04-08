//+build !faker

package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	ChecksParams struct {
		Owner      string `json:"owner" query:"owner"`
		Repository string `json:"repository" query:"repository"`
		Ref        string `json:"ref" query:"ref"`
	}
)

func (p *ChecksParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.Owner == "" {
		return &uiConfigModels.ConfigError{}
	}

	if p.Repository == "" {
		return &uiConfigModels.ConfigError{}
	}

	if p.Ref == "" {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}

// Used by cache as identifier
func (p *ChecksParams) String() string {
	return fmt.Sprintf("CHECKS-%s-%s-%s", p.Owner, p.Repository, p.Ref)
}
