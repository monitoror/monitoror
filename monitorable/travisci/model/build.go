package model

import "time"

type (
	Build struct {
		Branch string
		Author Author

		State         string //see https://github.com/shuheiktgw/go-travis/blob/master/builds.go#L133
		PreviousState string //see https://github.com/shuheiktgw/go-travis/blob/master/builds.go#L133
		StartedAt     time.Time
		FinishedAt    time.Time
		Duration      time.Duration //in second
	}

	Author struct {
		Name      string
		AvatarUrl string
	}
)
