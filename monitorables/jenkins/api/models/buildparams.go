//+build !faker

package models

import (
	"fmt"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
)

type (
	BuildParams struct {
		Job    string `json:"job" query:"job"`
		Branch string `json:"branch,omitempty" query:"branch,omitempty"`
	}
)

func (p *BuildParams) Validate(_ *uiConfigModels.ConfigVersion) *uiConfigModels.ConfigError {
	// TODO

	if p.Job == "" {
		return &uiConfigModels.ConfigError{}
	}

	return nil
}

// Used by cache as identifier
func (p *BuildParams) String() string {
	return fmt.Sprintf("BUILD-%s-%s", p.Job, p.Branch)
}
