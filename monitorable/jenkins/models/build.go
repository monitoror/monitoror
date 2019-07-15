package models

import "time"

type (
	Build struct {
		Number   string
		FullName string
		Author   *Author

		Building  bool
		Result    string
		StartedAt time.Time
		Duration  time.Duration
	}

	Author struct {
		Name      string
		AvatarUrl string
	}
)
