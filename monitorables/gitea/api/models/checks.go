package models

import "time"

type (
	Checks struct {
		HeadCommit *string
		Statuses   []Status
	}

	// Status represents a commit status reported by a CI (Gitea Actions, Drone, Woodpecker ...)
	// See : https://docs.gitea.com/api (CommitStatus)
	Status struct {
		ID        int64
		Title     string
		State     string // pending, success, error, failure, warning, skipped
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)
