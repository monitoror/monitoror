package models

import "time"

type (
	Release struct {
		ReleaseNumber  string
		DefinitionName string
		Author         *Author

		Status string

		FinishedAt *time.Time
		StartedAt  *time.Time
		QueuedAt   *time.Time
	}
)
