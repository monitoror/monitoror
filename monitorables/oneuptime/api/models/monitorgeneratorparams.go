package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	MonitorGeneratorParams struct {
		params.Default

		Labels string `json:"labels,omitempty" query:"labels"`
		SortBy string `json:"sortBy,omitempty" query:"sortBy" validate:"omitempty,oneof=name"`
	}
)
