//go:build faker

package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	"github.com/monitoror/monitoror/internal/pkg/validator"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	IssuesCountParams struct {
		params.Default

		Query                  string                  `json:"query" query:"query" validate:"required"`
		CountThreshold         coreModels.IntThreshold `json:"countThreshold" query:"countThreshold"`
		PriorityFieldThreshold coreModels.IntThreshold `json:"priorityFieldThreshold" query:"priorityFieldThreshold"`
		PriorityFieldLabel     *string                 `json:"priorityFieldLabel" query:"priorityFieldLabel"`
		PriorityFieldValue     *string                 `json:"priorityFieldValue" query:"priorityFieldValue"`

		Status      coreModels.TileStatus `json:"status" query:"status"`
		ValueValues []string              `json:"valueValues" query:"valueValues"`
	}
)

func (p *IssuesCountParams) GetPriorityQuery() string {
	return buildPriorityQuery(p.Query, p.PriorityFieldLabel, p.PriorityFieldValue)
}

func (p *IssuesCountParams) Validate() []validator.Error {
	return validateThresholds(p.CountThreshold, p.PriorityFieldThreshold)
}
