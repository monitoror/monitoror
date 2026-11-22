package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type PullRequestGeneratorParams struct {
	params.Default

	Owner      string `json:"owner" query:"owner" validate:"required"`
	Repository string `json:"repository" query:"repository" validate:"required"`
}
