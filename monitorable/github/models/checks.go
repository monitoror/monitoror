package models

import "time"

type (
	Checks struct {
		HeadCommit *string
		Runs       []Run
		Statuses   []Status
	}

	// Run represents GitHub Action Data
	// See : https://developer.github.com/v3/checks/runs/
	Run struct {
		ID          int64
		Title       string
		Status      string // queued, in_progress, or completed
		Conclusion  string // success, failure, neutral, cancelled, timed_out, or action_required
		StartedAt   *time.Time
		CompletedAt *time.Time
	}

	// Status Represent Simple "Status" provided by Apps (CodeCov, Travis ...)
	// See : https://developer.github.com/v3/repos/statuses/
	Status struct {
		ID        int64
		Title     string
		State     string // error, success, failure, or pending
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)
