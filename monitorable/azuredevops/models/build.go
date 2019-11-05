package models

import "time"

type (
	Build struct {
		BuildNumber    string
		DefinitionName string
		Branch         string
		Author         *Author

		Status string
		Result string

		FinishedAt *time.Time
		StartedAt  *time.Time
		QueuedAt   *time.Time
	}
)
