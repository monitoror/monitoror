package models

import "time"

type (
	Build struct {
		ID     uint
		Branch string
		Author Author

		State      string // see https://github.com/shuheiktgw/go-travis/blob/master/builds.go#L116
		StartedAt  time.Time
		FinishedAt time.Time
		Duration   time.Duration
	}

	Author struct {
		Name      string
		AvatarURL string
	}
)
