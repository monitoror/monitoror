package models

import "time"

type (
	Pipelines struct {
		HeadCommit string
		Runs       []Run
	}

	Run struct {
		ID         int
		Status     string
		Duration   int
		CreatedAt  time.Time
		StartedAt  *time.Time
		FinishedAt *time.Time
	}
)
