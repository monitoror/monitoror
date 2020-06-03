package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type MergeRequestGeneratorParams struct {
	params.Default

	ProjectID *int `json:"projectId" query:"projectId" validate:"required"`
}
