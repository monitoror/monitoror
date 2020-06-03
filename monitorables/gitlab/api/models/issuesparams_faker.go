//+build faker

package models

import "github.com/monitoror/monitoror/internal/pkg/monitorable/params"

type (
	IssuesParams struct {
		params.Default

		ProjectID *int `json:"projectId" query:"projectId"`

		State      *string  `json:"state" query:"state"`
		Labels     []string `json:"labels" query:"labels"`
		Milestone  *string  `json:"milestone" query:"milestone"`
		Scope      *string  `json:"scope" query:"scope"`
		Search     *string  `json:"search" query:"search"`
		AuthorID   *int     `json:"authorId" query:"authorId"`
		AssigneeID *int     `json:"assigneeId" query:"assigneeId"`

		ValueValues []string `json:"valueValues" query:"valueValues"`
	}
)
