//go:build !faker
// +build !faker

package models

import (
	"fmt"
	"slices"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	"github.com/monitoror/monitoror/internal/pkg/validator"
	coreModels "github.com/monitoror/monitoror/models"
)

const DefaultPriorityFieldLabel = "Priority"
const DefaultPriorityFieldValue = "Show-stopper"

type (
	IssuesCountParams struct {
		params.Default

		Query                  string                  `json:"query" query:"query" validate:"required"`
		CountThreshold         coreModels.IntThreshold `json:"countThreshold" query:"countThreshold"`
		PriorityFieldThreshold coreModels.IntThreshold `json:"priorityFieldThreshold" query:"priorityFieldThreshold"`
		PriorityFieldLabel     *string                 `json:"priorityFieldLabel" query:"PriorityFieldLabel"`
		PriorityFieldValue     *string                 `json:"priorityFieldValue" query:"PriorityFieldValue"`
	}
)

func (p *IssuesCountParams) GetPriorityFieldLabelWithDefault() string {
	result := DefaultPriorityFieldLabel
	if p.PriorityFieldLabel != nil {
		result = *p.PriorityFieldLabel
	}
	return result
}

func (p *IssuesCountParams) GetPriorityFieldValueWithDefault() string {
	result := DefaultPriorityFieldValue
	if p.PriorityFieldValue != nil {
		result = *p.PriorityFieldValue
	}
	return result
}

func (p *IssuesCountParams) Validate() []validator.Error {
	var errs []validator.Error

	if p.CountThreshold != nil && len(p.CountThreshold) > 0 {
		for tileStatus := range p.CountThreshold {
			if !slices.Contains(coreModels.ThresholdAuthorizedTileStatus, tileStatus) {
				errs = append(errs, validator.NewDefaultError(fmt.Sprintf("%s.%s", "CountThreshold", string(tileStatus)), fmt.Sprintf(" %v", coreModels.ThresholdAuthorizedTileStatus)))
			}
		}
	}

	if p.PriorityFieldThreshold != nil && len(p.PriorityFieldThreshold) > 0 {
		for tileStatus := range p.PriorityFieldThreshold {
			if !slices.Contains(coreModels.ThresholdAuthorizedTileStatus, tileStatus) {
				errs = append(errs, validator.NewDefaultError(fmt.Sprintf("%s.%s", "PriorityFieldThreshold", string(tileStatus)), fmt.Sprintf(" %v", coreModels.ThresholdAuthorizedTileStatus)))
			}
		}
	}

	return errs
}
