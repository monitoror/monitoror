package models

import "time"

type (
	Job struct {
		ID        string
		Buildable bool
		InQueue   bool
		QueuedAt  *time.Time

		Branches []string
	}
)
