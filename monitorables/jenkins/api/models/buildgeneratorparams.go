package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	BuildGeneratorParams struct {
		params.Default

		Job string `json:"job" query:"job" validate:"required"`

		// Using Match / Unmatch filter instead of one filter because Golang's standard regex library doesn't have negative look ahead.
		Match   string `json:"match,omitempty" query:"match" validate:"regex"`
		Unmatch string `json:"unmatch,omitempty" query:"unmatch" validate:"regex"`
	}
)
