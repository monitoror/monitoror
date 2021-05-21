//+build !faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
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
