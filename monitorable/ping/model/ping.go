package model

import "time"

type (
	Ping struct {
		Min     time.Duration
		Max     time.Duration
		Average time.Duration
	}
)
